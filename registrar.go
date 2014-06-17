package main

import (
  "log"
)

func Registrar(input chan []*FileEvent) {
  state := make(map[string]*FileState)
  for events := range input {
    log.Printf("Registrar received %d events\n", len(events))
    // Take the last event found for each file source
    for _, event := range events {
      // skip stdin
      if *event.Source == "-" {
        continue
      }
      // have to dereference the FileInfo here because os.FileInfo is an
      // interface, not a struct, so Go doesn't have smarts to call the Sys()
      // method on a pointer to os.FileInfo. :(
      ino, dev := file_ids(event.fileinfo)
      state[*event.Source] = &FileState{
        Source: event.Source,
        // take the offset + length of the line + newline char and
        // save it as the new starting offset.
        // This issues a problem, if the EOL is a CRLF! Then on start it read the LF again and generates a event with an empty line
        Offset: event.Offset + int64(len(*event.Text)) + 1,
        Inode:  ino,
        Device: dev,
      }
      //log.Printf("State %s: %d\n", *event.Source, event.Offset)
    }

    if len(state) > 0 {
      WriteRegistry(state, ".logstash-forwarder")
    }
  }
}
