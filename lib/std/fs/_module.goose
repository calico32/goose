native fn readFile(path)
native fn writeFile(path, data)
native fn appendFile(path, data)
native fn deleteFile(path)


export {
  readFile, writeFile, appendFile, deleteFile
}
