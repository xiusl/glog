<!DOCTYPE html>
<html>
<head>
  <title>Beego</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css" integrity="sha384-HSMxcRTRxnN+Bdg0JdbxYKrThecOKuH5zCYotlSAcp1+c8xmyTe9GYg1l9a69psu" crossorigin="anonymous">
  <style>
  </style>
</head>

<body>
<div class="container">
  <div>
    <nav class="navbar navbar-default">
      <div class="container-fluid">
        <div class="navbar-header">
          <a class="navbar-brand" href="">
            La
          </a>
        </div>
        <button id="add-button" type="button" class="btn btn-default navbar-btn">Add</button>
      </div>
    </nav>
  </div>
  <div>
    
    <table class="table table-bordered">
      <tr>
        <th>Key</th>
        <th>Config</th>
      </tr>

      {{range $idx, $kv := .Kvs}}
      <tr>
        <td>{{$kv.Key}}</td>
        <td>
          {{range $kv.Configs}}
            <kbd>{{.Topic}}</kbd> {{.Path}}
            <button class="setting-btn" data-topic="{{.Topic}}" data-path="{{.Path}}" data-key="{{$kv.Key}}" type="button" class="btn btn-default" aria-label="Left Align">
              <span class="glyphicon glyphicon-trash"></span>
            </button>
            <br/>
          {{end}}
        </td>
        
      </tr>
      {{end}}
    </table>
  </div>
</div>
<script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js" integrity="sha384-aJ21OjlMXNL5UyIl/XNwTMqvzeRMZH2w8c5cRVpzpU8Y5bApTppSuUkhZXN0VxHd" crossorigin="anonymous"></script>
<script>
$(document).ready(function(){
  $("#add-button").click(function(){
    window.location.href = "/add";
  });
  $(".setting-btn").click(function(e){
    console.log(e)
    var key = $(this).attr("data-key");
    var topic = $(this).attr("data-topic");
    var path = $(this).attr("data-path");
    console.log("key: ", key, "topic: ", topic, "path: ", path)
    const data = "key="+ key +"&path="+ path +"&topic="+ topic
    window.location.href = "/edit?"+ data
  });
});
</script>
</body>
</html>
