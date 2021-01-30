import Vue from 'vue';
import VueRouter from 'vue-router';
import axios from 'axios';
import dayjs from 'dayjs';
import localizedPlugin from 'dayjs/plugin/localizedFormat';
import 'bulma/css/bulma.css';
import { library } from '@fortawesome/fontawesome-svg-core';
import {
  faCircleNotch, faCog, faExclamationTriangle, faCheck, faEnvelope, faUser, faKey, faLock,
} from '@fortawesome/free-solid-svg-icons';
import {
  faFacebook, faGoogle,
} from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Home from './routes/home.vue';
import CreateAccount from './routes/createAccount.vue';
import LoginRedirect from './routes/loginRedirect.vue';
import ManageAccount from './routes/manageAccount.vue';
import PageNotFound from './routes/pageNotFound.vue';
import ForgotPassword from './routes/forgotPassword.vue';
import ActivateAccount from './routes/activateAccount.vue';
import OAuth2 from './routes/oauth2.vue';

axios.defaults.headers.common['X-CSRF-TOKEN'] = document.head.querySelector('meta[name="csrf"]').content;
dayjs.extend(localizedPlugin);

library.add(faExclamationTriangle, faCircleNotch, faCog, faCheck, faEnvelope, faUser, faKey, faLock);
library.add(faFacebook, faGoogle);
Vue.component('fa-icon', FontAwesomeIcon);

Vue.directive('focus', {
  inserted(ele) {
    ele.focus();
  },
});

window.bindRouter = function bindRouter(el, data = {}) {
  const router = new VueRouter({
    routes: [
      { path: '/', component: Home, props: data },
      { path: '/forgot-password', component: ForgotPassword, props: data },
      { path: '/create', component: CreateAccount, props: data },
      { path: '/login-redirect', component: LoginRedirect, props: data },
      { path: '/manage', component: ManageAccount, props: data },
      { path: '/activate', component: ActivateAccount, props: (route) => ({ token: route.query.token, account: route.query.account }) },
      {
        path: '/oauth2',
        component: OAuth2,
        props: (route) => ({
          meta: data,
          client_id: route.query.client_id,
          redirect_uri: route.query.redirect_uri,
          response_type: route.query.response_type,
          state: route.query.state,
          scope: route.query.scope,
        }),
      },
      { path: '*', component: PageNotFound },
    ],
  });
  Vue.use(VueRouter);

  return new Vue({
    el,
    router,
  });
};
