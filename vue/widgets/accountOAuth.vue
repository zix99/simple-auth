<template>
  <div>
    <LoadingBanner :promise="loadingPromise">Fetching Tokens...</LoadingBanner>
    <div v-for="token in this.tokens" :key="token.short_token" class="card my-2">
      <div class="card-content">
        <div class="media">
          <div class="media-left">
            <fa-icon icon="key" size="lg" />
          </div>
          <div class="media-content">
            <p class="title is-6">{{token.client_name}} ({{token.client_id}})</p>
            <div><strong>Key Type:</strong> <em>{{typeName(token.type)}}</em></div>
            <div><strong>Key Prefix:</strong> <code>{{token.short_token}}</code></div>
          </div>
        </div>
        <div class="content">
          <div class="light is-size-7">Added on <ShortDate :date="token.created" />, and expires on <ShortDate :date="token.expires" /></div>
        </div>
      </div>
    </div>
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
    typeName(t) {
      switch (t) {
        case 'access_token': return 'Access Token';
        case 'refresh_token': return 'Refresh Token';
        default: return 'Unknown';
      }
    },
  },
};
</script>
<style scoped>
div.light {
  color: #8a8a8a;
}
</style>
