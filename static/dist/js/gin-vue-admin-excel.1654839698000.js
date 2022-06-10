/*! 
 Build based on gin-vue-admin 
 Time : 1654839698000 */
import{s as e,i as a,_ as l,r as t,j as o,b as s,o as n,c as i,d,e as r,w as c,y as p,t as u,h as m,$ as h}from"../gva/gin-vue-admin-index.1654839698000.js";const v=(e,l)=>{if(void 0!==e.data){if("application/json"===e.data.type){const l=new FileReader;l.onload=function(){const e=JSON.parse(l.result).msg;a({showClose:!0,message:e,type:"error"})},l.readAsText(new Blob([e.data]))}}else{var t=window.URL.createObjectURL(new Blob([e])),o=document.createElement("a");o.style.display="none",o.href=t,o.download=l;var s=new MouseEvent("click");o.dispatchEvent(s)}},x=()=>e({url:"/excel/loadExcel",method:"get"});const f={class:"upload"},w={class:"gva-table-box"},b={class:"gva-btn-list"},g=m("导入"),y=m("导出"),E=m("下载模板");var _=l(Object.assign({name:"Excel"},{setup(a){const l=t("http://adminapi.seedreality.com"),m=t(1),_=t(0),k=t(999),j=t([]),z=(e=(()=>{}))=>{return a=this,l=null,t=function*(){const a=yield e({page:m.value,pageSize:k.value});0===a.code&&(j.value=a.data.list,_.value=a.data.total,m.value=a.data.page,k.value=a.data.pageSize)},new Promise(((e,o)=>{var s=e=>{try{i(t.next(e))}catch(a){o(a)}},n=e=>{try{i(t.throw(e))}catch(a){o(a)}},i=a=>a.done?e(a.value):Promise.resolve(a.value).then(s,n);i((t=t.apply(a,l)).next())}));var a,l,t};z(h);const I=o(),T=a=>{a&&"string"==typeof a||(a="ExcelExport.xlsx"),((a,l)=>{e({url:"/excel/exportExcel",method:"post",data:{fileName:l,infoList:a},responseType:"blob"}).then((e=>{v(e,l)}))})(j.value,a)},N=()=>{z(x)},C=()=>{var a;e({url:"/excel/downloadTemplate",method:"get",params:{fileName:a="ExcelTemplate.xlsx"},responseType:"blob"}).then((e=>{v(e,a)}))};return(e,a)=>{const t=s("el-button"),o=s("el-upload"),m=s("el-table-column"),h=s("el-table");return n(),i("div",f,[d("div",w,[d("div",b,[r(o,{class:"excel-btn",action:`${l.value}/excel/importExcel`,headers:{"x-token":p(I).token},"on-success":N,"show-file-list":!1},{default:c((()=>[r(t,{size:"small",type:"primary",icon:"upload"},{default:c((()=>[g])),_:1})])),_:1},8,["action","headers"]),r(t,{class:"excel-btn",size:"small",type:"primary",icon:"download",onClick:a[0]||(a[0]=e=>T("ExcelExport.xlsx"))},{default:c((()=>[y])),_:1}),r(t,{class:"excel-btn",size:"small",type:"success",icon:"download",onClick:a[1]||(a[1]=e=>C())},{default:c((()=>[E])),_:1})]),r(h,{data:j.value,"row-key":"ID"},{default:c((()=>[r(m,{align:"left",label:"ID","min-width":"100",prop:"ID"}),r(m,{align:"left","show-overflow-tooltip":"",label:"路由Name","min-width":"160",prop:"name"}),r(m,{align:"left","show-overflow-tooltip":"",label:"路由Path","min-width":"160",prop:"path"}),r(m,{align:"left",label:"是否隐藏","min-width":"100",prop:"hidden"},{default:c((e=>[d("span",null,u(e.row.hidden?"隐藏":"显示"),1)])),_:1}),r(m,{align:"left",label:"父节点","min-width":"90",prop:"parentId"}),r(m,{align:"left",label:"排序","min-width":"70",prop:"sort"}),r(m,{align:"left",label:"文件路径","min-width":"360",prop:"component"})])),_:1},8,["data"])])])}}}),[["__scopeId","data-v-4ff7d823"]]);export{_ as default};
