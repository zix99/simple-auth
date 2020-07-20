import Vue from 'vue';
import axios from 'axios';
import 'bulma/css/bulma.css';
import banner from './components/banner.vue';
import createaccount from './createAccount.vue';
import manageaccount from './manageAccount.vue';
import redirectinglogin from './redirectingLogin.vue';

axios.defaults.headers.common['X-CSRF-TOKEN'] = document.head.querySelector('meta[name="csrf"]').content;

window.bindVue = function bindVue(el, data = {}) {
  return new Vue({
    el,
    data,
    components: {
      banner,
      createaccount,
      manageaccount,
      redirectinglogin,
    },
  });
};
