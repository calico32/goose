import { ExtensionContext } from 'vscode'
import {
  Executable,
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  Trace,
} from 'vscode-languageclient/node'

export function activate(context: ExtensionContext) {
  const serverExecutable: Executable = {
    command: 'go',
    args: ['run', '/home/calico32/goose/bin/lsp/main.go'],
  }

  const serverOptions: ServerOptions = {
    run: serverExecutable,
    debug: serverExecutable,
  }

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ language: 'goose' }],
  }

  const client = new LanguageClient('Goose Language Server', serverOptions, clientOptions)
  client.setTrace(Trace.Verbose)
  client.start()

  context.subscriptions.push(client)
}
