<html>
    <head>
           <title>Upload file</title>
    </head>
    <body>
        <h1>Upload a sound file</h1>
    <form enctype="multipart/form-data" action="/putchime" method="post">
        <b>The sound file</b><input type="file" name="uploadfile" /><br/>
        <b>Optional path</b><input type="text" name="path" /><br/>
        <input type="hidden" name="token" value="{{.Token}}"/><br/>
        <input type="submit" value="upload" /><br/>
    </form>
    </body>
    </html>
