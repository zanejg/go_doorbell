<html>
    <head>
           <title>Upload file</title>
    </head>
    <body>
    <form enctype="multipart/form-data" action="http://{{.Serverstr}}:3400/putchime" method="post">
        <input type="file" name="uploadfile" />
        <input type="hidden" name="token" value="{{.Token}}"/>
        <input type="submit" value="upload" />
    </form>
    </body>
    </html>
