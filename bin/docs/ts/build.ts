import path from "path"

const root = path.resolve(import.meta.dir, "..")

await Bun.$`bun build ${root}/ts/highlight.ts --outfile=${root}/tmp/highlight.js`

// BUG: @shikijs/core/dist/wasm-inlined.mjs's init_wasm_inlined() is not called when
// shiki/dist/wasm.mjs is initialized. This leads to getWasmInstance being undefined
// and causing the error `wasm2.default is not a function`. To fix this, we need to
// replace:
//    var init_wasm2 = __esm(() => {
//    });
// with:
//    var init_wasm2 = __esm(() => {
//      init_wasm_inlined();
//    });
// in the generated highlight.js file.

let f = await Bun.file(`${root}/tmp/highlight.js`).text()
f = f.replace(
  "var init_wasm2 = __esm(() => {",
  "var init_wasm2 = __esm(() => {\n  init_wasm_inlined();"
)
await Bun.write(`${root}/tmp/highlight.js`, f)

await Bun.$`bun build --minify ${root}/tmp/highlight.js --outfile=${root}/tmp/highlight.js`
