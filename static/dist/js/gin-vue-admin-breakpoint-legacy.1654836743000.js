/*! 
 Build based on gin-vue-admin 
 Time : 1654836743000 */
!function(){function t(){"use strict";/*! regenerator-runtime -- Copyright (c) 2014-present, Facebook, Inc. -- license (MIT): https://github.com/facebook/regenerator/blob/main/LICENSE */t=function(){return e};var e={},r=Object.prototype,n=r.hasOwnProperty,a="function"==typeof Symbol?Symbol:{},i=a.iterator||"@@iterator",o=a.asyncIterator||"@@asyncIterator",u=a.toStringTag||"@@toStringTag";function f(t,e,r){return Object.defineProperty(t,e,{value:r,enumerable:!0,configurable:!0,writable:!0}),t[e]}try{f({},"")}catch(E){f=function(t,e,r){return t[e]=r}}function l(t,e,r,n){var a=e&&e.prototype instanceof h?e:h,i=Object.create(a.prototype),o=new k(n||[]);return i._invoke=function(t,e,r){var n="suspendedStart";return function(a,i){if("executing"===n)throw new Error("Generator is already running");if("completed"===n){if("throw"===a)throw i;return B()}for(r.method=a,r.arg=i;;){var o=r.delegate;if(o){var u=x(o,r);if(u){if(u===c)continue;return u}}if("next"===r.method)r.sent=r._sent=r.arg;else if("throw"===r.method){if("suspendedStart"===n)throw n="completed",r.arg;r.dispatchException(r.arg)}else"return"===r.method&&r.abrupt("return",r.arg);n="executing";var f=s(t,e,r);if("normal"===f.type){if(n=r.done?"completed":"suspendedYield",f.arg===c)continue;return{value:f.arg,done:r.done}}"throw"===f.type&&(n="completed",r.method="throw",r.arg=f.arg)}}}(t,r,o),i}function s(t,e,r){try{return{type:"normal",arg:t.call(e,r)}}catch(E){return{type:"throw",arg:E}}}e.wrap=l;var c={};function h(){}function p(){}function d(){}var v={};f(v,i,(function(){return this}));var y=Object.getPrototypeOf,g=y&&y(y(L([])));g&&g!==r&&n.call(g,i)&&(v=g);var m=d.prototype=h.prototype=Object.create(v);function b(t){["next","throw","return"].forEach((function(e){f(t,e,(function(t){return this._invoke(e,t)}))}))}function w(t,e){function r(a,i,o,u){var f=s(t[a],t,i);if("throw"!==f.type){var l=f.arg,c=l.value;return c&&"object"==typeof c&&n.call(c,"__await")?e.resolve(c.__await).then((function(t){r("next",t,o,u)}),(function(t){r("throw",t,o,u)})):e.resolve(c).then((function(t){l.value=t,o(l)}),(function(t){return r("throw",t,o,u)}))}u(f.arg)}var a;this._invoke=function(t,n){function i(){return new e((function(e,a){r(t,n,e,a)}))}return a=a?a.then(i,i):i()}}function x(t,e){var r=t.iterator[e.method];if(void 0===r){if(e.delegate=null,"throw"===e.method){if(t.iterator.return&&(e.method="return",e.arg=void 0,x(t,e),"throw"===e.method))return c;e.method="throw",e.arg=new TypeError("The iterator does not provide a 'throw' method")}return c}var n=s(r,t.iterator,e.arg);if("throw"===n.type)return e.method="throw",e.arg=n.arg,e.delegate=null,c;var a=n.arg;return a?a.done?(e[t.resultName]=a.value,e.next=t.nextLoc,"return"!==e.method&&(e.method="next",e.arg=void 0),e.delegate=null,c):a:(e.method="throw",e.arg=new TypeError("iterator result is not an object"),e.delegate=null,c)}function _(t){var e={tryLoc:t[0]};1 in t&&(e.catchLoc=t[1]),2 in t&&(e.finallyLoc=t[2],e.afterLoc=t[3]),this.tryEntries.push(e)}function A(t){var e=t.completion||{};e.type="normal",delete e.arg,t.completion=e}function k(t){this.tryEntries=[{tryLoc:"root"}],t.forEach(_,this),this.reset(!0)}function L(t){if(t){var e=t[i];if(e)return e.call(t);if("function"==typeof t.next)return t;if(!isNaN(t.length)){var r=-1,a=function e(){for(;++r<t.length;)if(n.call(t,r))return e.value=t[r],e.done=!1,e;return e.value=void 0,e.done=!0,e};return a.next=a}}return{next:B}}function B(){return{value:void 0,done:!0}}return p.prototype=d,f(m,"constructor",d),f(d,"constructor",p),p.displayName=f(d,u,"GeneratorFunction"),e.isGeneratorFunction=function(t){var e="function"==typeof t&&t.constructor;return!!e&&(e===p||"GeneratorFunction"===(e.displayName||e.name))},e.mark=function(t){return Object.setPrototypeOf?Object.setPrototypeOf(t,d):(t.__proto__=d,f(t,u,"GeneratorFunction")),t.prototype=Object.create(m),t},e.awrap=function(t){return{__await:t}},b(w.prototype),f(w.prototype,o,(function(){return this})),e.AsyncIterator=w,e.async=function(t,r,n,a,i){void 0===i&&(i=Promise);var o=new w(l(t,r,n,a),i);return e.isGeneratorFunction(r)?o:o.next().then((function(t){return t.done?t.value:o.next()}))},b(m),f(m,u,"Generator"),f(m,i,(function(){return this})),f(m,"toString",(function(){return"[object Generator]"})),e.keys=function(t){var e=[];for(var r in t)e.push(r);return e.reverse(),function r(){for(;e.length;){var n=e.pop();if(n in t)return r.value=n,r.done=!1,r}return r.done=!0,r}},e.values=L,k.prototype={constructor:k,reset:function(t){if(this.prev=0,this.next=0,this.sent=this._sent=void 0,this.done=!1,this.delegate=null,this.method="next",this.arg=void 0,this.tryEntries.forEach(A),!t)for(var e in this)"t"===e.charAt(0)&&n.call(this,e)&&!isNaN(+e.slice(1))&&(this[e]=void 0)},stop:function(){this.done=!0;var t=this.tryEntries[0].completion;if("throw"===t.type)throw t.arg;return this.rval},dispatchException:function(t){if(this.done)throw t;var e=this;function r(r,n){return o.type="throw",o.arg=t,e.next=r,n&&(e.method="next",e.arg=void 0),!!n}for(var a=this.tryEntries.length-1;a>=0;--a){var i=this.tryEntries[a],o=i.completion;if("root"===i.tryLoc)return r("end");if(i.tryLoc<=this.prev){var u=n.call(i,"catchLoc"),f=n.call(i,"finallyLoc");if(u&&f){if(this.prev<i.catchLoc)return r(i.catchLoc,!0);if(this.prev<i.finallyLoc)return r(i.finallyLoc)}else if(u){if(this.prev<i.catchLoc)return r(i.catchLoc,!0)}else{if(!f)throw new Error("try statement without catch or finally");if(this.prev<i.finallyLoc)return r(i.finallyLoc)}}}},abrupt:function(t,e){for(var r=this.tryEntries.length-1;r>=0;--r){var a=this.tryEntries[r];if(a.tryLoc<=this.prev&&n.call(a,"finallyLoc")&&this.prev<a.finallyLoc){var i=a;break}}i&&("break"===t||"continue"===t)&&i.tryLoc<=e&&e<=i.finallyLoc&&(i=null);var o=i?i.completion:{};return o.type=t,o.arg=e,i?(this.method="next",this.next=i.finallyLoc,c):this.complete(o)},complete:function(t,e){if("throw"===t.type)throw t.arg;return"break"===t.type||"continue"===t.type?this.next=t.arg:"return"===t.type?(this.rval=this.arg=t.arg,this.method="return",this.next="end"):"normal"===t.type&&e&&(this.next=e),c},finish:function(t){for(var e=this.tryEntries.length-1;e>=0;--e){var r=this.tryEntries[e];if(r.finallyLoc===t)return this.complete(r.completion,r.afterLoc),A(r),c}},catch:function(t){for(var e=this.tryEntries.length-1;e>=0;--e){var r=this.tryEntries[e];if(r.tryLoc===t){var n=r.completion;if("throw"===n.type){var a=n.arg;A(r)}return a}}throw new Error("illegal catch attempt")},delegateYield:function(t,e,r){return this.delegate={iterator:L(t),resultName:e,nextLoc:r},"next"===this.method&&(this.arg=void 0),c}},e}function e(t,e,r,n,a,i,o){try{var u=t[i](o),f=u.value}catch(l){return void r(l)}u.done?e(f):Promise.resolve(f).then(n,a)}function r(t){return function(){var r=this,n=arguments;return new Promise((function(a,i){var o=t.apply(r,n);function u(t){e(o,a,i,u,f,"next",t)}function f(t){e(o,a,i,u,f,"throw",t)}u(void 0)}))}}var n=document.createElement("style");n.innerHTML="h3[data-v-8564df28]{margin:40px 0 0}ul[data-v-8564df28]{list-style-type:none;padding:0}li[data-v-8564df28]{display:inline-block;margin:0 10px}a[data-v-8564df28]{color:#42b983}#fromCont[data-v-8564df28]{display:inline-block}.fileUpload[data-v-8564df28]{padding:3px 10px;font-size:12px;height:20px;line-height:20px;position:relative;cursor:pointer;color:#000;border:1px solid #c1c1c1;border-radius:4px;overflow:hidden;display:inline-block}.fileUpload input[data-v-8564df28]{position:absolute;font-size:100px;right:0;top:0;opacity:0;cursor:pointer}.fileName[data-v-8564df28]{display:inline-block;vertical-align:top;margin:6px 15px 0}.uploadBtn[data-v-8564df28]{position:relative;top:-10px;margin-left:15px}.tips[data-v-8564df28]{margin-top:30px;font-size:14px;font-weight:400;color:#606266}.el-divider[data-v-8564df28]{margin:0 0 30px}.list[data-v-8564df28]{margin-top:15px}.list-item[data-v-8564df28]{display:block;margin-right:10px;color:#606266;line-height:25px;margin-bottom:5px;width:40%}.list-item .percentage[data-v-8564df28]{float:right}.list-enter-active[data-v-8564df28],.list-leave-active[data-v-8564df28]{transition:all 1s}.list-enter[data-v-8564df28],.list-leave-to[data-v-8564df28]{opacity:0;transform:translateY(-30px)}\n",document.head.appendChild(n),System.register(["../gva/gin-vue-admin-index-legacy.1654836743000.js"],(function(e){"use strict";var n,a,i,o,u,f,l,s,c,h,p,d,v,y,g,m,b,w,x;return{setters:[function(t){n=t.s,a=t._,i=t.r,o=t.P,u=t.b,f=t.o,l=t.c,s=t.d,c=t.e,h=t.w,p=t.C,d=t.t,v=t.f,y=t.T,g=t.h,m=t.p,b=t.l,w=t.i,x=t.R}],execute:function(){var _={exports:{}};_.exports=function(t){var e=["0","1","2","3","4","5","6","7","8","9","a","b","c","d","e","f"];function r(t,e){var r=t[0],n=t[1],a=t[2],i=t[3];n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&a|~n&i)+e[0]-680876936|0)<<7|r>>>25)+n|0)&n|~r&a)+e[1]-389564586|0)<<12|i>>>20)+r|0)&r|~i&n)+e[2]+606105819|0)<<17|a>>>15)+i|0)&i|~a&r)+e[3]-1044525330|0)<<22|n>>>10)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&a|~n&i)+e[4]-176418897|0)<<7|r>>>25)+n|0)&n|~r&a)+e[5]+1200080426|0)<<12|i>>>20)+r|0)&r|~i&n)+e[6]-1473231341|0)<<17|a>>>15)+i|0)&i|~a&r)+e[7]-45705983|0)<<22|n>>>10)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&a|~n&i)+e[8]+1770035416|0)<<7|r>>>25)+n|0)&n|~r&a)+e[9]-1958414417|0)<<12|i>>>20)+r|0)&r|~i&n)+e[10]-42063|0)<<17|a>>>15)+i|0)&i|~a&r)+e[11]-1990404162|0)<<22|n>>>10)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&a|~n&i)+e[12]+1804603682|0)<<7|r>>>25)+n|0)&n|~r&a)+e[13]-40341101|0)<<12|i>>>20)+r|0)&r|~i&n)+e[14]-1502002290|0)<<17|a>>>15)+i|0)&i|~a&r)+e[15]+1236535329|0)<<22|n>>>10)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&i|a&~i)+e[1]-165796510|0)<<5|r>>>27)+n|0)&a|n&~a)+e[6]-1069501632|0)<<9|i>>>23)+r|0)&n|r&~n)+e[11]+643717713|0)<<14|a>>>18)+i|0)&r|i&~r)+e[0]-373897302|0)<<20|n>>>12)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&i|a&~i)+e[5]-701558691|0)<<5|r>>>27)+n|0)&a|n&~a)+e[10]+38016083|0)<<9|i>>>23)+r|0)&n|r&~n)+e[15]-660478335|0)<<14|a>>>18)+i|0)&r|i&~r)+e[4]-405537848|0)<<20|n>>>12)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&i|a&~i)+e[9]+568446438|0)<<5|r>>>27)+n|0)&a|n&~a)+e[14]-1019803690|0)<<9|i>>>23)+r|0)&n|r&~n)+e[3]-187363961|0)<<14|a>>>18)+i|0)&r|i&~r)+e[8]+1163531501|0)<<20|n>>>12)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n&i|a&~i)+e[13]-1444681467|0)<<5|r>>>27)+n|0)&a|n&~a)+e[2]-51403784|0)<<9|i>>>23)+r|0)&n|r&~n)+e[7]+1735328473|0)<<14|a>>>18)+i|0)&r|i&~r)+e[12]-1926607734|0)<<20|n>>>12)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n^a^i)+e[5]-378558|0)<<4|r>>>28)+n|0)^n^a)+e[8]-2022574463|0)<<11|i>>>21)+r|0)^r^n)+e[11]+1839030562|0)<<16|a>>>16)+i|0)^i^r)+e[14]-35309556|0)<<23|n>>>9)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n^a^i)+e[1]-1530992060|0)<<4|r>>>28)+n|0)^n^a)+e[4]+1272893353|0)<<11|i>>>21)+r|0)^r^n)+e[7]-155497632|0)<<16|a>>>16)+i|0)^i^r)+e[10]-1094730640|0)<<23|n>>>9)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n^a^i)+e[13]+681279174|0)<<4|r>>>28)+n|0)^n^a)+e[0]-358537222|0)<<11|i>>>21)+r|0)^r^n)+e[3]-722521979|0)<<16|a>>>16)+i|0)^i^r)+e[6]+76029189|0)<<23|n>>>9)+a|0,n=((n+=((a=((a+=((i=((i+=((r=((r+=(n^a^i)+e[9]-640364487|0)<<4|r>>>28)+n|0)^n^a)+e[12]-421815835|0)<<11|i>>>21)+r|0)^r^n)+e[15]+530742520|0)<<16|a>>>16)+i|0)^i^r)+e[2]-995338651|0)<<23|n>>>9)+a|0,n=((n+=((i=((i+=(n^((r=((r+=(a^(n|~i))+e[0]-198630844|0)<<6|r>>>26)+n|0)|~a))+e[7]+1126891415|0)<<10|i>>>22)+r|0)^((a=((a+=(r^(i|~n))+e[14]-1416354905|0)<<15|a>>>17)+i|0)|~r))+e[5]-57434055|0)<<21|n>>>11)+a|0,n=((n+=((i=((i+=(n^((r=((r+=(a^(n|~i))+e[12]+1700485571|0)<<6|r>>>26)+n|0)|~a))+e[3]-1894986606|0)<<10|i>>>22)+r|0)^((a=((a+=(r^(i|~n))+e[10]-1051523|0)<<15|a>>>17)+i|0)|~r))+e[1]-2054922799|0)<<21|n>>>11)+a|0,n=((n+=((i=((i+=(n^((r=((r+=(a^(n|~i))+e[8]+1873313359|0)<<6|r>>>26)+n|0)|~a))+e[15]-30611744|0)<<10|i>>>22)+r|0)^((a=((a+=(r^(i|~n))+e[6]-1560198380|0)<<15|a>>>17)+i|0)|~r))+e[13]+1309151649|0)<<21|n>>>11)+a|0,n=((n+=((i=((i+=(n^((r=((r+=(a^(n|~i))+e[4]-145523070|0)<<6|r>>>26)+n|0)|~a))+e[11]-1120210379|0)<<10|i>>>22)+r|0)^((a=((a+=(r^(i|~n))+e[2]+718787259|0)<<15|a>>>17)+i|0)|~r))+e[9]-343485551|0)<<21|n>>>11)+a|0,t[0]=r+t[0]|0,t[1]=n+t[1]|0,t[2]=a+t[2]|0,t[3]=i+t[3]|0}function n(t){var e,r=[];for(e=0;e<64;e+=4)r[e>>2]=t.charCodeAt(e)+(t.charCodeAt(e+1)<<8)+(t.charCodeAt(e+2)<<16)+(t.charCodeAt(e+3)<<24);return r}function a(t){var e,r=[];for(e=0;e<64;e+=4)r[e>>2]=t[e]+(t[e+1]<<8)+(t[e+2]<<16)+(t[e+3]<<24);return r}function i(t){var e,a,i,o,u,f,l=t.length,s=[1732584193,-271733879,-1732584194,271733878];for(e=64;e<=l;e+=64)r(s,n(t.substring(e-64,e)));for(a=(t=t.substring(e-64)).length,i=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],e=0;e<a;e+=1)i[e>>2]|=t.charCodeAt(e)<<(e%4<<3);if(i[e>>2]|=128<<(e%4<<3),e>55)for(r(s,i),e=0;e<16;e+=1)i[e]=0;return o=(o=8*l).toString(16).match(/(.*?)(.{0,8})$/),u=parseInt(o[2],16),f=parseInt(o[1],16)||0,i[14]=u,i[15]=f,r(s,i),s}function o(t){var e,n,i,o,u,f,l=t.length,s=[1732584193,-271733879,-1732584194,271733878];for(e=64;e<=l;e+=64)r(s,a(t.subarray(e-64,e)));for(n=(t=e-64<l?t.subarray(e-64):new Uint8Array(0)).length,i=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],e=0;e<n;e+=1)i[e>>2]|=t[e]<<(e%4<<3);if(i[e>>2]|=128<<(e%4<<3),e>55)for(r(s,i),e=0;e<16;e+=1)i[e]=0;return o=(o=8*l).toString(16).match(/(.*?)(.{0,8})$/),u=parseInt(o[2],16),f=parseInt(o[1],16)||0,i[14]=u,i[15]=f,r(s,i),s}function u(t){var r,n="";for(r=0;r<4;r+=1)n+=e[t>>8*r+4&15]+e[t>>8*r&15];return n}function f(t){var e;for(e=0;e<t.length;e+=1)t[e]=u(t[e]);return t.join("")}function l(t){return/[\u0080-\uFFFF]/.test(t)&&(t=unescape(encodeURIComponent(t))),t}function s(t,e){var r,n=t.length,a=new ArrayBuffer(n),i=new Uint8Array(a);for(r=0;r<n;r+=1)i[r]=t.charCodeAt(r);return e?i:a}function c(t){return String.fromCharCode.apply(null,new Uint8Array(t))}function h(t,e,r){var n=new Uint8Array(t.byteLength+e.byteLength);return n.set(new Uint8Array(t)),n.set(new Uint8Array(e),t.byteLength),r?n:n.buffer}function p(t){var e,r=[],n=t.length;for(e=0;e<n-1;e+=2)r.push(parseInt(t.substr(e,2),16));return String.fromCharCode.apply(String,r)}function d(){this.reset()}return f(i("hello")),"undefined"==typeof ArrayBuffer||ArrayBuffer.prototype.slice||function(){function e(t,e){return(t=0|t||0)<0?Math.max(t+e,0):Math.min(t,e)}ArrayBuffer.prototype.slice=function(r,n){var a,i,o,u,f=this.byteLength,l=e(r,f),s=f;return n!==t&&(s=e(n,f)),l>s?new ArrayBuffer(0):(a=s-l,i=new ArrayBuffer(a),o=new Uint8Array(i),u=new Uint8Array(this,l,a),o.set(u),i)}}(),d.prototype.append=function(t){return this.appendBinary(l(t)),this},d.prototype.appendBinary=function(t){this._buff+=t,this._length+=t.length;var e,a=this._buff.length;for(e=64;e<=a;e+=64)r(this._hash,n(this._buff.substring(e-64,e)));return this._buff=this._buff.substring(e-64),this},d.prototype.end=function(t){var e,r,n=this._buff,a=n.length,i=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];for(e=0;e<a;e+=1)i[e>>2]|=n.charCodeAt(e)<<(e%4<<3);return this._finish(i,a),r=f(this._hash),t&&(r=p(r)),this.reset(),r},d.prototype.reset=function(){return this._buff="",this._length=0,this._hash=[1732584193,-271733879,-1732584194,271733878],this},d.prototype.getState=function(){return{buff:this._buff,length:this._length,hash:this._hash.slice()}},d.prototype.setState=function(t){return this._buff=t.buff,this._length=t.length,this._hash=t.hash,this},d.prototype.destroy=function(){delete this._hash,delete this._buff,delete this._length},d.prototype._finish=function(t,e){var n,a,i,o=e;if(t[o>>2]|=128<<(o%4<<3),o>55)for(r(this._hash,t),o=0;o<16;o+=1)t[o]=0;n=(n=8*this._length).toString(16).match(/(.*?)(.{0,8})$/),a=parseInt(n[2],16),i=parseInt(n[1],16)||0,t[14]=a,t[15]=i,r(this._hash,t)},d.hash=function(t,e){return d.hashBinary(l(t),e)},d.hashBinary=function(t,e){var r=f(i(t));return e?p(r):r},d.ArrayBuffer=function(){this.reset()},d.ArrayBuffer.prototype.append=function(t){var e,n=h(this._buff.buffer,t,!0),i=n.length;for(this._length+=t.byteLength,e=64;e<=i;e+=64)r(this._hash,a(n.subarray(e-64,e)));return this._buff=e-64<i?new Uint8Array(n.buffer.slice(e-64)):new Uint8Array(0),this},d.ArrayBuffer.prototype.end=function(t){var e,r,n=this._buff,a=n.length,i=[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];for(e=0;e<a;e+=1)i[e>>2]|=n[e]<<(e%4<<3);return this._finish(i,a),r=f(this._hash),t&&(r=p(r)),this.reset(),r},d.ArrayBuffer.prototype.reset=function(){return this._buff=new Uint8Array(0),this._length=0,this._hash=[1732584193,-271733879,-1732584194,271733878],this},d.ArrayBuffer.prototype.getState=function(){var t=d.prototype.getState.call(this);return t.buff=c(t.buff),t},d.ArrayBuffer.prototype.setState=function(t){return t.buff=s(t.buff,!0),d.prototype.setState.call(this,t)},d.ArrayBuffer.prototype.destroy=d.prototype.destroy,d.ArrayBuffer.prototype._finish=d.prototype._finish,d.ArrayBuffer.hash=function(t,e){var r=f(o(new Uint8Array(t)));return e?p(r):r},d}();var A=_.exports,k=function(t){return n({url:"/fileUploadAndDownload/findFile",method:"get",params:t})},L=function(t){return n({url:"/fileUploadAndDownload/breakpointContinueFinish",method:"post",params:t})},B=function(t,e){return n({url:"/fileUploadAndDownload/removeChunk",method:"post",data:t,params:e})},E=function(t){return m("data-v-8564df28"),t=t(),b(),t},C={class:"break-point"},S={class:"gva-table-box"},U=g("大文件上传"),F={id:"fromCont",method:"post"},N=g(" 选择文件 "),j=g("上传文件"),I=E((function(){return s("div",{class:"el-upload__tip"},"请上传不超过5MB的文件",-1)})),M={class:"list"},O={key:0,class:"list-item"},P={class:"percentage"},D=E((function(){return s("div",{class:"tips"},"此版本为先行体验功能测试版，样式美化和性能优化正在进行中，上传切片文件和合成的完整文件分别再QMPlusserver目录的breakpointDir文件夹和fileDir文件夹",-1)})),T={name:"BreakPoint"},G=Object.assign(T,{setup:function(e){var a=i(null),g=i(""),m=i([]),b=i([]),_=i(NaN),E=i(!1),T=i(0),G=i(!0),z=function(){var e=r(t().mark((function e(n){var i,o;return t().wrap((function(e){for(;;)switch(e.prev=e.next){case 0:i=new FileReader,o=n.target.files[0],5242880,a.value=o,T.value=0,a.value.size<5242880?(i.readAsArrayBuffer(a.value),i.onload=function(){var e=r(t().mark((function e(r){var n,i,o,u,f,l,s,c,h,p,d;return t().wrap((function(t){for(;;)switch(t.prev=t.next){case 0:for(n=r.target.result,(i=new A.ArrayBuffer).append(n),g.value=i.end(),o=1048576,u=0,f=0,l=0,m.value=[];f<a.value.size;)u=l*o,f=(l+1)*o,s=a.value.slice(u,f),(c=new window.FormData).append("fileMd5",g.value),c.append("file",s),c.append("chunkNumber",l),c.append("fileName",a.value.name),m.value.push({key:l,formData:c}),l++;return h={fileName:a.value.name,fileMd5:g.value,chunkTotal:m.value.length},t.next=13,k(h);case 13:p=t.sent,d=p.data.file.ExaFileChunk,p.data.file.IsFinish?(b.value=[],w.success("文件已秒传")):b.value=m.value.filter((function(t){return!(d&&d.some((function(e){return e.FileChunkNumber===t.key})))})),_.value=b.value.length,console.log(_.value);case 19:case"end":return t.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}()):(E.value=!0,w("请上传小于5M文件"));case 6:case"end":return e.stop()}}),e)})));return function(t){return e.apply(this,arguments)}}(),R=function(){null!==a.value?(100===T.value&&(G.value=!1),Y()):w("请先上传文件")},Y=function(){b.value&&b.value.forEach((function(t){t.formData.append("chunkTotal",m.value.length);var e=new FileReader,r=t.formData.get("file");e.readAsArrayBuffer(r),e.onload=function(e){var r=new A.ArrayBuffer;r.append(e.target.result),t.formData.append("chunkMd5",r.end()),$(t)}}))};o((function(){return _.value}),(function(){T.value=Math.floor((m.value.length-_.value)/m.value.length*100)}));var $=function(){var e=r(t().mark((function e(r){var i,o,u;return t().wrap((function(t){for(;;)switch(t.prev=t.next){case 0:return t.next=2,e=r.formData,n({url:"/fileUploadAndDownload/breakpointContinue",method:"post",headers:{"Content-Type":"multipart/form-data"},data:e});case 2:if(0===t.sent.code){t.next=5;break}return t.abrupt("return");case 5:if(_.value--,0!==_.value){t.next=16;break}return i={fileName:a.value.name,fileMd5:g.value},t.next=10,L(i);case 10:if(0!==(o=t.sent).code){t.next=16;break}return u={fileName:a.value.name,fileMd5:g.value,filePath:o.data.filePath},w.success("上传成功"),t.next=16,B(u);case 16:case"end":return t.stop()}var e}),e)})));return function(t){return e.apply(this,arguments)}}(),H=i(null),Q=function(){H.value.dispatchEvent(new MouseEvent("click"))};return function(t,e){var r=u("el-divider"),n=u("el-button"),i=u("document"),o=u("el-icon"),g=u("el-progress");return f(),l("div",C,[s("div",S,[c(r,{"content-position":"left"},{default:h((function(){return[U]})),_:1}),s("form",F,[s("div",{class:"fileUpload",onClick:Q},[N,p(s("input",{id:"file",ref_key:"FileInput",ref:H,multiple:"multiple",type:"file",onChange:z},null,544),[[x,!1]])])]),c(n,{disabled:E.value,type:"primary",size:"small",class:"uploadBtn",onClick:R},{default:h((function(){return[j]})),_:1},8,["disabled"]),I,s("div",M,[c(y,{name:"list",tag:"p"},{default:h((function(){return[a.value?(f(),l("div",O,[c(o,null,{default:h((function(){return[c(i)]})),_:1}),s("span",null,d(a.value.name),1),s("span",P,d(T.value)+"%",1),c(g,{"show-text":!1,"text-inside":!1,"stroke-width":2,percentage:T.value},null,8,["percentage"])])):v("",!0)]})),_:1})]),D])])}}});e("default",a(G,[["__scopeId","data-v-8564df28"]]))}}}))}();
