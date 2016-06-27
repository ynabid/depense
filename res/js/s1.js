window.onload = function(){
  var today=new Date();
  var day=today.getDate();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  if(month<10) month ="0"+month;
  if(day<10) day="0"+day;
  loadAll();
  document.getElementById("date").value = year+"-"+month+"-"+day;
  document.getElementById("login_form").onsubmit = function(){
    var email = this.querySelector("[name=email]").value;
    var password = this.querySelector("[name=password]").value;
    login(email,password,function(response){
      loadAll();
      document.getElementById("login_bg").style.display = "none"
    });
    return false;	
  } 
  document.getElementById("depense").onsubmit = function(){

    var xhr;
    if (window.XMLHttpRequest) {
      xhr = new XMLHttpRequest();
    } else {
      // code for IE6, IE5
      xhr = new ActiveXObject("Microsoft.XMLHTTP");
    }
    xhr.open("POST","/api/depense",true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function(){
      if (xhr.readyState == 4){
	switch(xhr.status){
	  case 201 :
	    loadAllMonth()
	      document.getElementById("description").value="";
	    document.getElementById("montant").value="";

	    break;
	  default :
	    console.log(xhr.response);
	}
      }
    }
    var data = JSON.stringify({
      date : document.getElementById("date").value,
      description : document.getElementById("description").value,
      montant : document.getElementById("montant").value,
      categoryid : document.getElementById("category").value

    });

    xhr.send(data);
    return false;	
  } 
}
function getCurrentMonth(){
  var today=new Date();
  var day=today.getDate();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  if(month<10) month ="0"+month;
  return year+"-"+month;
}
function getXHR(){
  var xhr;
  if (window.XMLHttpRequest) {
    xhr = new XMLHttpRequest();
  } else {
    // code for IE6, IE5
    xhr = new ActiveXObject("Microsoft.XMLHTTP");
  }
  return xhr;
}
function loadAll(){
  loadCategories();
  loadAllMonth(); 
  loadDepenseByCategory();
}
function loadDepenseByCategory(){
  var xhr=getXHR();
  var currentMonth = getCurrentMonth();
  xhr.open("GET","/api/depense/bycategory?month="+currentMonth,true);
  xhr.onreadystatechange = function(){
    if (this.readyState == 4){
      switch(this.status){
	case 200 :
	  var dCatList = JSON.parse(xhr.responseText);
	  var p=document.getElementById("listCat");
	  var le= p.firstElementChild.cloneNode(true);
	  while (p.firstChild) {
	    p.removeChild(p.firstChild);
	  }
	  for(d in dCatList){
	    var lei= le.cloneNode(true);
	    lei.querySelector("[name=id]").innerHTML = dCatList[d].Category.Id;
	    lei.querySelector("[name=category]").innerHTML = dCatList[d].Category.Name;
	    lei.querySelector("[name=montant]").innerHTML = dCatList[d].Montant;
	    p.appendChild(lei)
	  }
	  break;
	default :
	  console.log(this.response);
      }
    }
  }
  xhr.send();
}

function loadCategories(){
  var xhr;
  if (window.XMLHttpRequest) {
    xhr = new XMLHttpRequest();
  } else {
    // code for IE6, IE5
    xhr = new ActiveXObject("Microsoft.XMLHTTP");
  }

  xhr.open("GET","/api/depense/category/list",true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function(){
    if (this.readyState == 4){
      switch(this.status){
	case 200 :
	  var categories = JSON.parse(xhr.responseText);
	  var categoryLV = document.getElementById("category");
	  for (var i = 0, len = categories.length; i < len; i++) {
	    var category = categories[i];
	    var option = document.createElement("option");
	    option.value = category.Id;
	    if (category.Id == 1){
	      option.selected = true;
	    }
	    option.textContent = category.Name;
	    categoryLV.add(option);
	  }
	  break;
	default :
	  console.log(this.response);
      }
    }
  }
  xhr.send();
}
function loadAllMonth(){
  if(isAuth()){
    loadPPMonth();
    loadThisMonth();
    loadPMonth();
    loadDepenseListTM();
  }
}
function loadThisMonth(){
  var today=new Date();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  loadData(year,month,function(res){
    document.getElementById("this-month").innerHTML = res;
  })
}
function loadPMonth(){
  var today=new Date();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  if(month == 1){
    year--;
    month = 12;
  }else{
    month--	
  }
  loadData(year,month,function(res){
    document.getElementById("p-month").innerHTML = res;
  })
}

function loadPPMonth(){
  var today=new Date();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  if(month <= 2){
    year--;
    month = 13 - month;
  }else{
    month= month -2;	
  }
  loadData(year,month,function(res){
    document.getElementById("pp-month").innerHTML = res;
  })
}

function loadData(year,month,callback){
  if(month<10) month ="0"+month;
  var xhr;
  if (window.XMLHttpRequest) {
    xhr = new XMLHttpRequest();
  } else {
    // code for IE6, IE5
    xhr = new ActiveXObject("Microsoft.XMLHTTP");
  }
  xhr.open("GET","/api/depense/month"+"?month="+year+"-"+month,true);
  xhr.onreadystatechange = function(){
    if (xhr.readyState == 4){
      switch(xhr.status){
	case 200 :
	  callback(xhr.responseText)
	    break;
	default :
	  console.log(xhr.response);
      }
    }
  }
  xhr.send();
}
function loadDepenseListTM(){
  var today=new Date();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  loadDepenseList(
      year,
      month,
      function(dList){
	var p=document.getElementById("list");
	var le= p.firstElementChild.cloneNode(true);
	while (p.firstChild) {
	  p.removeChild(p.firstChild);
	}
	for(d in dList){
	  var lei= le.cloneNode(true);
	  lei.querySelector("[name=id]").innerHTML = dList[d].Id;
	  lei.querySelector("[name=category]").innerHTML = dList[d].Category;
	  lei.querySelector("[name=description]").innerHTML = dList[d].Description;
	  lei.querySelector("[name=montant]").innerHTML = dList[d].Montant;
	  var date = new Date(dList[d].Date*1000);
	  var day = date.getDate();
	  var m = ""
	    if(day < 10) 
	      m="0"+day;
	    else
	      m= day;
	  lei.querySelector("[name=date]").innerHTML = m; 
	  p.appendChild(lei)
	}
      }
  );
}
function loadDepenseList(year,month,callback){
  if(month<10) month ="0"+month;
  var xhr;
  if (window.XMLHttpRequest) {
    xhr = new XMLHttpRequest();
  } else {
    // code for IE6, IE5
    xhr = new ActiveXObject("Microsoft.XMLHTTP");
  }
  xhr.open("GET","/api/depenseList"+"?month="+year+"-"+month,true);
  xhr.onreadystatechange = function(){
    if (xhr.readyState == 4){
      switch(xhr.status){
	case 200 :
	  callback(JSON.parse(xhr.responseText));
	  break;
	default :
	  console.log(xhr.response);
      }
    }
  }
  xhr.send();
}
function getCookie(cname) {
  var name = cname + "=";
  var ca = document.cookie.split(';');
  for(var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}
function checkCookie(name) {
  var coockie=getCookie(name);
  if (coockie!="") {
    return true;
  } else {
    return false
  }

}
function isAuth(){
  if(!checkCookie("uid")){	
    document.getElementById("login_bg").style.display = "block"
      var loginPopup = document.getElementById("login");
    loginPopup.style.left = (window.screen.width - loginPopup.offsetWidth)/2 + "px";
    loginPopup.style.top = (window.screen.height - loginPopup.offsetHeight)/2 + "px";
    return false;
  }
  return true;
}
function login(email,password,callback){
  var xhr;
  if (window.XMLHttpRequest) {
    xhr = new XMLHttpRequest();
  } else {
    // code for IE6, IE5
    xhr = new ActiveXObject("Microsoft.XMLHTTP");
  }
  xhr.open("POST","/api/auth",true);
  xhr.onreadystatechange = function(){
    if (xhr.readyState == 4){
      switch(xhr.status){
	case 200 :
	  callback(xhr.response);
	  break;
	default :
	  console.log(xhr.response);
      }
    }
  }
  xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
  xhr.send("email="+email+"&password="+password);
}
