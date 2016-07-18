var depenses;
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
    addEntry(
	{
	  date : document.getElementById("date").value,
	  description : document.getElementById("description").value,
	  //type : Number(type),
	  montant : Number(document.getElementById("montant").value),
	  categoryid : Number(document.getElementById("category").value),
	  accountid : Number(document.getElementById("account").value)
	},function(){
	  loadAll();
	  document.getElementById("description").value="";
	  document.getElementById("montant").value="";
	}
	);
    return false;	
  } 
}

function addEntry(d,callback){
  var xhr=getXHR();
  xhr.open("POST","/api/depense",true);
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onreadystatechange = function(){
    if (xhr.readyState == 4){
      switch(xhr.status){
	case 201 :
	  callback();
	  break;
	default :
	  console.log(xhr.response);
      }
    }
  }
  var data = JSON.stringify(d);
  xhr.send(data);
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
  var xhr=getXHR();
  var currentMonth = getCurrentMonth();
  xhr.open("GET","/api/depense/all?month="+currentMonth,true);
  xhr.onreadystatechange = function(){
    if (this.readyState == 4){
      switch(this.status){
	case 200 :
	  depenses = JSON.parse(xhr.responseText);
	  bindAccountsTR();
	  bindDepenseByCategory();
	  bindDepenseList();
	  bindAccounts();
	  bindCategories();

      }
    }
  }
  xhr.send();
}
function bindAccountsTR(){
  var row=document.getElementById("listAccountHead");
  var col= row.firstElementChild.cloneNode(true);

  var p2=document.getElementById("listAccount");
  var row2= p2.firstElementChild.cloneNode(true);
  var col2= row2.firstElementChild.cloneNode(true);



  while (row.firstChild) {
    row.removeChild(row.firstChild);
  }

  while (p2.firstChild) {
    p2.removeChild(p2.firstChild);
  }
  while (row2.firstChild) {
    row2.removeChild(row2.firstChild);
  }


  var accountsTR = depenses.AccountsTR;
  var accounts = depenses.Accounts;
  accounts.push({Id : -1, Name :""});
  var ch = col.cloneNode(true);
  ch.innerHTML = "";
  row.appendChild(ch);

  for(i in accounts){
    var ch = col.cloneNode(true);
    ch.innerHTML = accounts[i].Name;
    row.appendChild(ch);
  }

  for(i in accounts){
    var r = row2.cloneNode(true);
    var c = col2.cloneNode(true);
    c.innerHTML = accounts[i].Name;
    r.appendChild(c);
    for(j in accounts){
      c = col2.cloneNode(true);
      if(accountsTR[accounts[i].Name] && accountsTR[accounts[i].Name][accounts[j].Name]){
	c.innerHTML =  accountsTR[accounts[i].Name][accounts[j].Name]; 
      }else{
	c.innerHTML = 0;
      }
      r.appendChild(c);
    }
    p2.appendChild(r);
  }
}

function bindDepenseByCategory(){
  var p=document.getElementById("listCatHead");
  var le= p.firstElementChild.cloneNode(true);
  var p2=document.getElementById("catTotal");
  var le2= p2.firstElementChild.cloneNode(true);
  var accounts = depenses.Accounts;
  var categories = depenses.Categories;
  var dCatList = depenses.DepenseCategory;

  while (p.firstChild) {
    p.removeChild(p.firstChild);
  }
  p.appendChild(le);

  while (p2.firstChild) {
    p2.removeChild(p2.firstChild);
  }
  p2.appendChild(le2);


  var lei= le.cloneNode(true);
  lei.innerHTML = "Total";
  p.appendChild(lei);

  var lei2= le2.cloneNode(true);
  if(dCatList["Total"]["Total"]){
    lei2.innerHTML = dCatList["Total"]["Total"];
  }else{
    lei2.innerHTML = 0; 
  }
  p2.appendChild(lei2);


  for(i in accounts){
    var account = accounts[i];
    var lei= le.cloneNode(true);
    lei.innerHTML = account.Name;
    p.appendChild(lei);

    var lei2= le2.cloneNode(true);
    if(dCatList["Total"][account.Name]){
      lei2.innerHTML = dCatList["Total"][account.Name];
    }else{
      lei2.innerHTML = 0; 
    }
    p2.appendChild(lei2);

  }

  var p=document.getElementById("listCat");
  var le= p.firstElementChild.cloneNode(true);
  var lei = le.firstElementChild;
  while (p.firstChild) {
    p.removeChild(p.firstChild);
  }
  while (le.firstChild) {
    le.removeChild(le.firstChild);
  }
  le.appendChild(lei);

  for(i in categories){
    var category = categories[i];
    var lei= le.cloneNode(true);
    var cell = lei.firstElementChild;
    cell.innerHTML = category.Name;
    lei.appendChild(cell);

    var cell2 = cell.cloneNode(true);
    if(dCatList[category.Name] && dCatList[category.Name]["Total"]){
      cell2.innerHTML = dCatList[category.Name]["Total"];
    }else{
      cell2.innerHTML = 0;
    }
    lei.appendChild(cell2);


    for(j in accounts){
      var cell2 = cell.cloneNode(true);
      var account = accounts[j];
      if(dCatList[category.Name] && dCatList[category.Name][account.Name]){
	cell2.innerHTML = dCatList[category.Name][account.Name];
      }else{
	cell2.innerHTML = 0;
      }
      lei.appendChild(cell2);
    }
    p.appendChild(lei)
  }
}
function bindAccounts(){
  accounts = depenses.Accounts;
  var accountLV = document.getElementById("account");
  while (accountLV.firstChild) {
    accountLV.removeChild(accountLV.firstChild);
  }

  for (var i = 0, len = accounts.length; i < len; i++) {
    var account = accounts[i];
    var option = document.createElement("option");
    option.value = account.Id;
    if (account.Id == 1){
      option.selected = true;
    }
    option.textContent = account.Name;
    accountLV.add(option);
  }
}

function bindCategories(){
  categories = depenses.Categories;
  var categoryLV = document.getElementById("category");
  while (categoryLV.firstChild) {
    categoryLV.removeChild(categoryLV.firstChild);

  }

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
}
function loadData(year,month,callback){
  var xhr=getXHR();
  if(month<10) month ="0"+month;
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
function convertEpoch(t){
  var date = new Date(t*1000);
  var day = date.getDate();
  var month=date.getMonth()+1;
  var year=date.getFullYear();
  if(month<10) month ="0"+month;
  if(day < 10) 
    day="0"+day;
  return year+"-"+month+"-"+day;
}

function bindDepenseList(){
  var row=document.getElementById("listDepenseHead");
  var col= row.firstElementChild.cloneNode(true);

  var p2=document.getElementById("list");
  var row2= p2.firstElementChild.cloneNode(true);
  var col2= row2.firstElementChild.cloneNode(true);



  while (row.firstChild) {
    row.removeChild(row.firstChild);
  }

  while (p2.firstChild) {
    p2.removeChild(p2.firstChild);
  }
  while (row2.firstChild) {
    row2.removeChild(row2.firstChild);
  }


  var depensesList = depenses.Depenses;
  var depenseHead = depensesList[0];
  if(depenseHead){
    for(colName in depenseHead){
      var ch = col.cloneNode(true);
      ch.innerHTML = colName;
      row.appendChild(ch);
    }

    for(i in depensesList){
      var r = row2.cloneNode(true);
      var c = col2.cloneNode(true);
      var dl=depensesList[i];
      for(j in dl){
	var dc = dl[j];
	c = col2.cloneNode(true);
	if(j=="Date"){
	  var date = new Date(dc);
	  var day = date.getDate();
	  if(day < 10) 
	    day="0"+day;
	  c.innerHTML =  day; 
	}else{
	  c.innerHTML =  dc; 
	}
	r.appendChild(c);
      }
      p2.appendChild(r);
    }
  }
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
    document.getElementById("login_bg").style.display = "block";
    var loginPopup = document.getElementById("login");
    loginPopup.style.left = (window.screen.width - loginPopup.offsetWidth)/2 + "px";
    loginPopup.style.top = (window.screen.height - loginPopup.offsetHeight)/2 + "px";
    return false;
  }
  return true;
}
function login(email,password,callback){
  var xhr=getXHR();
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
