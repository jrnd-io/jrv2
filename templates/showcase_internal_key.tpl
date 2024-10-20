{{ $id:=uuid}}{{.SetKey $id}}{
  "key": "{{$id}}"
}{{.AddHeader "key" $id}}{{.AddHeader "time" (now "2006-01-02T15:04:05Z07:00")}}
