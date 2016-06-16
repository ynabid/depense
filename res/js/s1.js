window.onload = function(){
  var today=new Date();
  var day=today.getDate();
  var month=today.getMonth()+1;
  var year=today.getFullYear();
  if(month<10) month ="0"+month;
  if(day<10) day="0"+day;
  loadAllMonth();
  document.getElementById("date").value = year+"-"+month+"-"+day;
  document.getElementById("depense").onsubmit = function(){

    var xhr;
    if (window.XMLHttpRequest) {
      xhr = new XMLHttpRequest();
    } else {
      // code for IE6, IE5
      xhr = new ActiveXObject("Microsoft.XMLHTTP");
    }
    xhr.open("POST","/api/depense/",true);
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
      montant : document.getElementById("montant").value
    });

    xhr.send(data);
    return false;	
  } 
}
function loadAllMonth(){
  loadPPMonth();
  loadThisMonth();
  loadPMonth();
  loadDepenseListTM();
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
	  lei.querySelector("[name=description]").innerHTML = dList[d].Description;
	  lei.querySelector("[name=montant]").innerHTML = dList[d].Montant;
	  var date = new Date(dList[d].Date*1000);
	  var month = date.getMonth()+1;
	  var m = ""
	    if(month < 10) 
	      m="0"+month;
	    else
	      m= month;
	  lei.querySelector("[name=date]").innerHTML = date.getFullYear()+"-"+m+"-"+date.getDate(); 
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
