
/**
 * First we will load all of this project's JavaScript dependencies which
 * includes Vue and other libraries. It is a great starting point when
 * building robust, powerful web applications using Vue and Laravel.
 */

require('./bootstrap');

window.Vue = require('vue');

// Basic Use - Covers most scenarios

import Croppa from 'vue-croppa';

import BlockUI from 'vue-blockui';

Vue.use(BlockUI);

Vue.use(Croppa);

Vue.component("Desk",require("./components/Desk.vue"))


/**
 * Next, we will create a fresh Vue application instance and attach it to
 * the page. Then, you may begin adding components to this application
 * or customize the JavaScript scaffolding to fit your unique needs.
 */

const app = new Vue({
    el: '#app'
});
