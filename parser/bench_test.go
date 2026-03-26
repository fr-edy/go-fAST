package parser_test

import (
	"strings"
	"testing"

	"github.com/t14raptor/go-fast/generator"
	"github.com/t14raptor/go-fast/parser"
)

var smallJS = `var x = 1 + 2 * 3;`

var mediumJS = `
function fibonacci(n) {
  if (n <= 1) return n;
  let a = 0, b = 1;
  for (let i = 2; i <= n; i++) {
    let tmp = b;
    b = a + b;
    a = tmp;
  }
  return b;
}

class EventEmitter {
  constructor() {
    this.listeners = {};
  }
  on(event, cb) {
    if (!this.listeners[event]) this.listeners[event] = [];
    this.listeners[event].push(cb);
    return this;
  }
  emit(event, ...args) {
    const cbs = this.listeners[event];
    if (cbs) for (const cb of cbs) cb(...args);
  }
  off(event, cb) {
    const cbs = this.listeners[event];
    if (cbs) this.listeners[event] = cbs.filter(f => f !== cb);
  }
}

const obj = {
  a: 1,
  b: "hello",
  c: [1, 2, 3],
  d: { nested: true },
  e: function() { return this.a; },
  f: () => 42,
  get g() { return this.b.length; },
  set g(v) { this.b = String(v); },
  [Symbol.iterator]() { return this.c[Symbol.iterator](); },
};

async function fetchData(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) throw new Error("HTTP " + response.status);
    const data = await response.json();
    return data.items.map(item => ({
      id: item.id,
      name: item.name,
      tags: [...item.tags],
    }));
  } catch (err) {
    console.error("Failed:", err.message);
    return [];
  }
}

function* range(start, end, step = 1) {
  for (let i = start; i < end; i += step) {
    yield i;
  }
}

const template = ` + "`Hello ${name}, you have ${count} items.`" + `;

const destructured = ({ a, b: { c, d = 10 }, ...rest }) => [a, c, d, rest];

switch (action.type) {
  case "INCREMENT": return { ...state, count: state.count + 1 };
  case "DECREMENT": return { ...state, count: state.count - 1 };
  case "RESET": return { ...state, count: 0 };
  default: return state;
}
`

var largeJS string

func init() {
	var b strings.Builder
	for i := 0; i < 50; i++ {
		b.WriteString(mediumJS)
		b.WriteByte('\n')
	}
	largeJS = b.String()
}

func BenchmarkParseSmall(b *testing.B) {
	for b.Loop() {
		_, _ = parser.ParseFile(smallJS)
	}
}

func BenchmarkParseMedium(b *testing.B) {
	for b.Loop() {
		_, _ = parser.ParseFile(mediumJS)
	}
}

func BenchmarkParseLarge(b *testing.B) {
	for b.Loop() {
		_, _ = parser.ParseFile(largeJS)
	}
}

func BenchmarkRoundTripMedium(b *testing.B) {
	for b.Loop() {
		p, _ := parser.ParseFile(mediumJS)
		_ = generator.Generate(p)
	}
}

func BenchmarkRoundTripLarge(b *testing.B) {
	for b.Loop() {
		p, _ := parser.ParseFile(largeJS)
		_ = generator.Generate(p)
	}
}

func BenchmarkGenerateMedium(b *testing.B) {
	p, _ := parser.ParseFile(mediumJS)
	b.ResetTimer()
	for b.Loop() {
		_ = generator.Generate(p)
	}
}

func BenchmarkGenerateLarge(b *testing.B) {
	p, _ := parser.ParseFile(largeJS)
	b.ResetTimer()
	for b.Loop() {
		_ = generator.Generate(p)
	}
}
