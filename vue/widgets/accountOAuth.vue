<template>
  <div>
    <LoadingBanner :promise="loadingPromise">Fetching Tokens...</LoadingBanner>
    <div class="card my-2" v-for="(tokenGroup, clientId) in this.groupedTokens" :key="clientId">
      <div class="card-content">
        <p class="title is-6">{{tokenGroup[0].client_name}} ({{clientId}})</p>
        <div class="media" v-for="token in tokenGroup" :key="token.short_token">
          <div class="media-left">
            <fa-icon icon="key" size="lg" />
          </div>
          <div class="media-content">
            <div><strong>Key Type:</strong> <em>{{typeName(token.type)}}</em></div>
            <div><strong>Key Prefix:</strong> <code>{{token.short_token}}</code></div>
            <div class="light is-size-7">Added on <ShortDate :date="token.created" />, and expires on <ShortDate :date="token.expires" /></div>
          </div>
        </div>
        <div class="buttons is-right">
          <button class="button is-danger is-small is-outlined" @click="revokeClientTokens(clientId)">Revoke {{tokenGroup[0].client_name}}'s Tokens</button>
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
  computed: {
    groupedTokens() {
      if (!this.tokens) return {};
      return this.tokens.reduce((memo, val) => {
        const k = val.client_id;
        (memo[k] = memo[k] || []).push(val); // eslint-disable-line no-param-reassign
        return memo;
      }, {});
    },
  },
  methods: {
    fetchData() {
      this.loadingPromise = axios.get('api/v1/auth/oauth2')
        .then((resp) => {
          this.tokens = resp.data.tokens;
        });
    },
    revokeClientTokens(clientId) {
      this.loadingPromise = axios.delete('api/v1/auth/oauth2/token', { params: { client_id: clientId } })
        .then(() => this.fetchData());
    },
    typeName(t) {
      switch (t) {
        case 'access_token': return 'Access Token';
        case 'refresh_token': return 'Refresh Token';
        case 'code': return 'Temporary Code';
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
