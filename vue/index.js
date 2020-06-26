import Vue from 'vue'
import Banner from './components/banner.vue'

window.bindVue = function bindVue(el, data = {}) {
    new Vue({
        el,
        data,
        components: {
            Banner,
        },
    });
}
