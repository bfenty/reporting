{{define "loginjs"}}
$(document).ready(function () {
  $("form").submit(function (event) {
    var formData = {
      username: $("#username").val(),
      password: $("#password").val(),
    };
    $.ajax({
      type: "POST",
      url: "/signin",
      data: formData,
      dataType: "application/json",
      encode: false,
      success: location.Reload(),
    }).done(function (data) {
      console.log(data);
    });
    event.preventDefault();
  });
});
{{end}}

{{define "skujs"}}
$(document).ready(function () {
  $("form").submit(function (event) {
    var formData = {
      sku_internal: $("#sku_internal").val(),
    };
    $.ajax({
      type: "POST",
      url: "/skus",
      data: formData,
      dataType: "application/json",
      encode: false,
      success: location.Reload(),
    }).done(function (data) {
      console.log(data);
    });
    event.preventDefault();
  });
});
{{end}}
