<template>
  <div>
    <LoadingBanner :promise="loadingPromise">Fetching Tokens...</LoadingBanner>
    <table class="table is-striped" v-if="tokens">
      <thead>
        <tr>
          <th>Client</th>
          <th>Token</th>
          <th>Issued</th>
          <th>Expires</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="token in this.tokens" :key="token.short_token">
          <td>{{token.client_name}}</td>
          <td>{{token.short_token}}</td>
          <td><ShortDate :date="token.created" /></td>
          <td><ShortDate :date="token.expires" /></td>
          <td></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from 'axios';
import LoadingBanner from '../components/loadingBanner.vue';
import ShortDate from '../components/shortdate.vue';

export default {
  components: {
    LoadingBanner,
    ShortDate,
  },
  data() {
    return {
      loadingPromise: null,
      tokens: [],
    };
  },
  created() {
    this.fetchData();
  },
  methods: {
    fetchData() {
      this.loadingPromise = axios.get('api/v1/auth/oauth2')
        .then((resp) => {
          this.tokens = resp.data.tokens;
        });
    },
  },
};
</script>
