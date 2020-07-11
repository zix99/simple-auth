import Vue from 'vue';
import banner from './components/banner.vue';
import createaccount from './createAccount.vue';
import manageaccount from './manageAccount.vue';

window.bindVue = function bindVue(el, data = {}) {
  return new Vue({
    el,
    data,
    components: {
      banner,
      createaccount,
      manageaccount,
    },
  });
};
