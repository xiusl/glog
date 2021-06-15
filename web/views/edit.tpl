<!DOCTYPE html>
<html>
<head>
  <title>LogAgent - Edit</title>
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
          <a class="navbar-brand" href="/">
            La - Edit
          </a>
        </div>
      </div>
    </nav>
  </div>
  <div>
    <div>
        {{.flash.error}}
        {{.flash.warning}}
        {{.flash.notice}}
    </div>
    <form style="width:320px;margin:60px auto;" action="/edit" method="post">
      <div class="form-group">
        <label for="key">Key</label>
        <input type="text" class="form-control" name="key" id="key" placeholder="Key">
      </div>
      <div class="form-group">
        <label for="topic">Topic</label>
        <input type="text" class="form-control" name="topic" id="topic" placeholder="Topic">
      </div>
      <div class="form-group">
        <label for="path">Path</label>
        <input type="text" class="form-control" name="path" id="path" placeholder="Path">
      </div>
      <button type="submit" class="btn btn-danger">Delete</button>
    </form>
  </div>
</div>
<script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.js"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js" integrity="sha384-aJ21OjlMXNL5UyIl/XNwTMqvzeRMZH2w8c5cRVpzpU8Y5bApTppSuUkhZXN0VxHd" crossorigin="anonymous"></script>
<script>
function getUrlParam(name) {
  var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
  var r = window.location.search.substr(1).match(reg);
  if (r != null) {
    return unescape(r[2]);
  }
  return null; 
}
var okey = getUrlParam("key")
var opath = getUrlParam("path")
var otopic = getUrlParam("topic")
$(document).ready(function(){
    $("#key").val(okey);
    $("#path").val(opath);
    $("#topic").val(otopic);
});
</script>
</body>
</html>
