<html>
    <head>
           <title>List MP3 files</title>
    </head>
    <body>
        <h1>List of Chimes on Doorbell Central</h1>
        {{$svr := .Ringserver}}
        <ul>
            {{ range $filename := .Filelist}}
            <li>
                {{$filename}}: <a href = "{{$svr}}{{$filename}}">ring</a>
            </li>
            {{end}}
        </ul>


    </body>
</html>
