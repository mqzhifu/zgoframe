/*! 
 Build based on gin-vue-admin 
 Time : 1654836743000 */
import{s as a}from"../gva/gin-vue-admin-index.1654836743000.js";const t=t=>a({url:"/autoCode/preview",method:"post",data:t}),o=t=>a({url:"/autoCode/createTemp",method:"post",data:t,responseType:"blob"}),e=()=>a({url:"/autoCode/getDB",method:"get"}),d=t=>a({url:"/autoCode/getTables",method:"get",params:t}),s=t=>a({url:"/autoCode/getColumn",method:"get",params:t}),u=t=>a({url:"/autoCode/getSysHistory",method:"post",data:t}),r=t=>a({url:"/autoCode/rollback",method:"post",data:t}),l=t=>a({url:"/autoCode/getMeta",method:"post",data:t}),m=t=>a({url:"/autoCode/delSysHistory",method:"post",data:t}),p=t=>a({url:"/autoCode/createPackage",method:"post",data:t}),g=()=>a({url:"/autoCode/getPackage",method:"post"}),h=t=>a({url:"/autoCode/delPackage",method:"post",data:t}),C=t=>a({url:"/autoCode/createPlug",method:"post",data:t});export{s as a,g as b,o as c,e as d,l as e,u as f,d as g,m as h,p as i,h as j,C as k,t as p,r};
