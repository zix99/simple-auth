import Vue from 'vue';
import VueRouter from 'vue-router';
import axios from 'axios';
import dayjs from 'dayjs';
import localizedPlugin from 'dayjs/plugin/localizedFormat';
import 'bulma/css/bulma.css';
import Home from './routes/home.vue';
import CreateAccount from './routes/createAccount.vue';
import LoginRedirect from './routes/loginRedirect.vue';
import ManageAccount from './routes/manageAccount.vue';
import PageNotFound from './routes/pageNotFound.vue';
import ForgotPassword from './routes/forgotPassword.vue';

axios.defaults.headers.common['X-CSRF-TOKEN'] = document.head.querySelector('meta[name="csrf"]').content;
dayjs.extend(localizedPlugin);

window.bindRouter = function bindRouter(el, data = {}) {
  const router = new VueRouter({
    routes: [
      { path: '/', component: Home, props: data },
      { path: '/forgot-password', component: ForgotPassword, props: data },
      { path: '/create', component: CreateAccount, props: data },
      { path: '/login-redirect', component: LoginRedirect, props: data },
      { path: '/manage', component: ManageAccount, props: data },
      { path: '*', component: PageNotFound },
    ],
  });
  Vue.use(VueRouter);
  return new Vue({
    el,
    router,
  });
};
