/*! 
 Build based on gin-vue-admin 
 Time : 1654839698000 */
import{_ as e,r as a,j as r,N as s,b as c,o as i,c as p,F as t,y as l,q as n,f as u}from"../gva/gin-vue-admin-index.1654839698000.js";const d={class:"headerAvatar"},o=["src"],m=["src"],g=["src"],y={name:"CustomPic"};var v=e(Object.assign(y,{props:{picType:{type:String,required:!1,default:"avatar"},picSrc:{type:String,required:!1,default:""}},setup(e){const y=e,v=a("http://adminapi.seedreality.com/"),I=a("./assets/noBody.745c3d16.png"),f=r(),S=s((()=>""===y.picSrc?""!==f.userInfo.headerImg&&"http"===f.userInfo.headerImg.slice(0,4)?f.userInfo.headerImg:v.value+f.userInfo.headerImg:""!==y.picSrc&&"http"===y.picSrc.slice(0,4)?y.picSrc:v.value+y.picSrc)),h=s((()=>y.picSrc&&"http"!==y.picSrc.slice(0,4)?v.value+y.picSrc:y.picSrc));return(a,r)=>{const s=c("el-avatar");return i(),p("span",d,["avatar"===e.picType?(i(),p(t,{key:0},[l(f).userInfo.headerImg?(i(),n(s,{key:0,size:30,src:l(S)},null,8,["src"])):(i(),n(s,{key:1,size:30,src:I.value},null,8,["src"]))],64)):u("",!0),"img"===e.picType?(i(),p(t,{key:1},[l(f).userInfo.headerImg?(i(),p("img",{key:0,src:l(S),class:"avatar"},null,8,o)):(i(),p("img",{key:1,src:I.value,class:"avatar"},null,8,m))],64)):u("",!0),"file"===e.picType?(i(),p("img",{key:2,src:l(h),class:"file"},null,8,g)):u("",!0)])}}}),[["__scopeId","data-v-47cae38c"]]);export{v as C};
