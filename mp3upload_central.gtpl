<html>
    <head>
           <title>Upload file</title>
    </head>
    <body>
        <h1>Upload a chime</h1>

        <h2>{{.Reply}}</h2>
    <form enctype="multipart/form-data" action="http://{{.Serverstr}}:3434/getchime" method="post">
        <input type="file" name="uploadfile" />
        <input type="hidden" name="token" value="{{.Token}}"/>
        <input type="submit" value="upload" />
    </form>
    </body>
    </html>
