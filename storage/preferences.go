/*
 * Copyright 2021 Aletheia Ware LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"aletheiaware.com/bcgo"
	"aletheiaware.com/spaceclientgo"
	"aletheiaware.com/spacego"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type preferences struct {
	ctx       context.Context
	client    spaceclientgo.SpaceClient
	node      bcgo.Node
	name      string
	values    map[string]string
	lock      sync.RWMutex
	listeners []func()
}

func NewPreferences(ctx context.Context, client spaceclientgo.SpaceClient, node bcgo.Node, name string) fyne.Preferences {
	p := &preferences{
		ctx:    ctx,
		client: client,
		node:   node,
		name:   name,
		values: make(map[string]string),
	}
	if metaId := p.metaId(); len(metaId) > 0 {
		p.watch(metaId)
	}
	return p
}

func (p *preferences) Bool(key string) bool {
	return p.BoolWithFallback(key, false)
}

func (p *preferences) BoolWithFallback(key string, fallback bool) bool {
	b, err := strconv.ParseBool(p.StringWithFallback(key, strconv.FormatBool(fallback)))
	if err != nil {
		return fallback
	}
	return b
}

func (p *preferences) SetBool(key string, value bool) {
	p.SetString(key, strconv.FormatBool(value))
}

func (p *preferences) Float(key string) float64 {
	return p.FloatWithFallback(key, 0.0)
}

func (p *preferences) FloatWithFallback(key string, fallback float64) float64 {
	f, err := strconv.ParseFloat(p.StringWithFallback(key, strconv.FormatFloat(fallback, 'E', -1, 64)), 64)
	if err != nil {
		return fallback
	}
	return f
}

func (p *preferences) SetFloat(key string, value float64) {
	p.SetString(key, strconv.FormatFloat(value, 'E', -1, 64))
}

func (p *preferences) Int(key string) int {
	return p.IntWithFallback(key, 0)
}

func (p *preferences) IntWithFallback(key string, fallback int) int {
	i, err := strconv.ParseInt(p.StringWithFallback(key, strconv.FormatInt(int64(fallback), 10)), 10, 64)
	if err != nil {
		return fallback
	}
	return int(i)
}

func (p *preferences) SetInt(key string, value int) {
	p.SetString(key, strconv.FormatInt(int64(value), 10))
}

func (p *preferences) String(key string) string {
	return p.StringWithFallback(key, "")
}

func (p *preferences) StringWithFallback(key, fallback string) string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	v, ok := p.values[key]
	if !ok {
		return fallback
	}
	return v
}

func (p *preferences) SetString(key string, value string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	v, ok := p.values[key]
	if ok && value == v {
		// No change
		return
	}
	p.values[key] = value

	p.write()
}

func (p *preferences) RemoveValue(key string) {
	delete(p.values, key)
	p.write()
}

func (p *preferences) AddChangeListener(listener func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.listeners = append(p.listeners, listener)
}

func (p *preferences) metaId() []byte {
	var (
		metaId    []byte
		timestamp uint64
	)
	// Choose the most recent preference file with matching name
	if err := p.client.SearchMeta(p.node, spacego.NewNameFilter(p.name), func(e *bcgo.BlockEntry, m *spacego.Meta) error {
		if t := e.Record.Timestamp; t > timestamp {
			timestamp = t
			metaId = e.RecordHash
		}
		return nil
	}); err != nil {
		fyne.LogError("Failed to find preferences file", err)
		return nil
	}
	return metaId
}

func (p *preferences) watch(metaId []byte) {
	p.client.WatchFile(p.ctx, p.node, metaId, func() {
		reader, err := p.client.ReadFile(p.node, metaId)
		if err != nil {
			fyne.LogError("Failed to read preferences file", err)
			return
		}

		scanner := bufio.NewScanner(reader)
		p.lock.Lock()
		for scanner.Scan() {
			parts := strings.SplitN(strings.TrimSpace(scanner.Text()), "=", 2)
			key := parts[0]
			var value string
			if len(parts) > 1 {
				value = parts[1]
			}
			p.values[key] = value
		}
		p.lock.Unlock()

		if err := scanner.Err(); err != nil {
			fyne.LogError("Failed to parse preferences file", err)
			return
		}

		// Trigger Listeners
		p.lock.RLock()
		for _, l := range p.listeners {
			go l()
		}
		p.lock.RUnlock()
	})
}

func (p *preferences) write() {
	var keys []string
	for k := range p.values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, p.values[k]))
	}
	data := []byte(sb.String())

	if metaId := p.metaId(); len(metaId) == 0 {
		// Create new preference file
		ref, err := p.client.Add(p.node, nil, p.name, spacego.MIME_TYPE_TEXT_PLAIN, bytes.NewReader(data))
		if err != nil {
			fyne.LogError("Failed to create preferences file", err)
			return
		}
		p.watch(ref.RecordHash)
	} else {
		// Amend existing preference file
		writer, err := p.client.WriteFile(p.node, nil, metaId)
		if err != nil {
			fyne.LogError("Failed to write preferences file", err)
			return
		}
		count, err := writer.Write(data)
		if err != nil {
			fyne.LogError("Failed to write preferences file", err)
			return
		}
		if count != len(data) {
			fyne.LogError(fmt.Sprintf("Failed to write all preferences: Expected %d, Wrote %d", len(data), count), nil)
			return
		}
		if err := writer.Close(); err != nil {
			fyne.LogError("Failed to write preferences file", err)
			return
		}
	}
}
