import { transformerNotationErrorLevel, transformerNotationHighlight } from "@shikijs/transformers"
import * as shiki from "shiki"
import grammar from "../../../extension/syntaxes/goose.tmLanguage.json"
import theme from "../theme.json"

async function highlight(code: string) {
  const highligher = await shiki.getHighlighter({
    langs: [
      {
        ...(grammar as any),
        name: "goose",
        scopeName: "source.goose",
      },
    ],
    themes: [],
  })

  const html = highligher.codeToHtml(code, {
    lang: "goose",
    theme: theme as any,
    transformers: [transformerNotationHighlight(), transformerNotationErrorLevel()],
  })

  return html
}

console.log(await highlight(await Bun.stdin.text()))
