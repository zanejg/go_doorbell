<html>
    <head>
           <title>Upload file</title>
    </head>
    <body>
        <h1>Upload a chime</h1>

        <h2>{{.Reply}}</h2>
    <form enctype="multipart/form-data" action="/getchime" method="post">
        <input type="file" name="uploadfile" />
        <input type="hidden" name="token" value="{{.Token}}"/>
        <input type="submit" value="upload" />
    </form>

    <a href = "/RingAllDoorbells">Ring all doorbells</a>
    </body>
    </html>
