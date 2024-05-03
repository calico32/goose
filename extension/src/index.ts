import { ExtensionContext } from "vscode"
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  Trace,
  TransportKind,
} from "vscode-languageclient/node"

export function activate(context: ExtensionContext) {
  let serverExecutable = {
    command: context.extensionPath + "/dist/goose-lsp",
    transport: TransportKind.stdio,
  }

  try {
    const config = require("./config.json")
    serverExecutable = config
  } catch {
    // ignore
  }

  const serverOptions: ServerOptions = {
    run: serverExecutable,
    debug: serverExecutable,
  }

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ language: "goose" }],
  }

  const client = new LanguageClient(
    "gooseLanguageServer",
    "Goose Language Server",
    serverOptions,
    clientOptions
  )
  client.setTrace(Trace.Verbose)
  client.start()

  context.subscriptions.push(client)
}
