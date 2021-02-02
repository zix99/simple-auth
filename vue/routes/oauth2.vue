<template>
  <CenterCard title="OAuth Login">
    <h2 class="subtitle">{{meta.appdata.company}} Login</h2>

    <LoadingBanner :promise="loadingPromise">
      <template v-slot:default>
        Granting access...
      </template>
      <template v-slot:success>
        Success! Redirecting...
      </template>
    </LoadingBanner>

    <p>The following application is requesting access to login to your account.</p>

    <div v-if="error">
      <Message type="is-danger">
        {{error}}
      </Message>
    </div>
    <div v-else-if="fetchingMeta">
      <Message type="is-info">
        <fa-icon icon="circle-notch" spin /> Fetching user information...
      </Message>
    </div>
    <div v-else>
      <h3>{{client.name || `Client: ${client_id}`}}</h3>
      <strong>By:</strong> <a target="_blank" :href="client.author_url">{{client.author}} ({{client.author_url}})</a>

      <div v-if="account" class="box">
        <strong>You are currently logged in as:</strong><br />
        {{account.email}}
        <div class="my-2">
          <button class="button is-small is-outlined" @click="switchAccountClick">Not you? Switch Account</button>
        </div>
      </div>
      <div v-if="noaccount" class="box">
        <p>You are not currently logged in</p>
      </div>
      <div v-if="scope">
        <strong>Requested Permissions:</strong>
        {{scope.replace(' ', ', ')}}
      </div>

      <div class="buttons is-right">
        <button class="button is-light" @click="cancelClick">Cancel</button>
        <button class="button is-primary" :class="{ 'is-loading': fetchingGrant }" :disabled="!account || fetchingGrant" @click="grantClick">Grant</button>
      </div>
    </div>
  </CenterCard>
</template>

<script>
import axios from 'axios';
import CenterCard from '../components/centerCard.vue';
import Message from '../components/message.vue';
import LoadingBanner from '../components/loadingBanner.vue';

const errorCodes = {
  invalid_client: 'Unknown client',
  invalid_scope: 'Invalid scopes',
};

export default {
  props: {
    meta: {},

    // From params
    client_id: null,
    response_type: null,
    redirect_uri: null,
    state: null,
    scope: null,
  },
  data() {
    return {
      error: null,
      loadingPromise: null,
      fetchingMeta: true,
      fetchingGrant: true,
      account: null,
      noaccount: false,
      client: {},
    };
  },
  components: {
    CenterCard,
    Message,
    LoadingBanner,
  },
  created() {
    if (!this.client_id) {
      this.error = 'Missing client id';
      return;
    }
    if (!this.redirect_uri) {
      this.error = 'Missing redirect URI';
      return;
    }
    if (this.response_type !== 'code') {
      this.error = 'Expected response_type of code';
      return;
    }

    axios.get('api/v1/account')
      .then((resp) => {
        this.account = resp.data;

        // Attempt to auto-grant (if the token previously exists)
        // Only do this if we have a session (otherwise it's moot)
        this.grant(true)
          .catch(() => { this.fetchingMeta = false; });
      }).catch(() => {
        this.noaccount = true;
        this.redirectToLogin();
      });

    axios.get(`api/v1/auth/oauth2/client/${this.client_id}`)
      .then((resp) => {
        this.client = resp.data;
      }).catch((err) => {
        this.error = errorCodes[err.response.data.error];
      });
  },
  methods: {
    cancelClick() {
      window.location = this.client.author_url;
    },
    grantClick() {
      this.loadingPromise = this.grant();
    },
    switchAccountClick() {
      this.account = null;
      this.noaccount = true;
      axios.delete('api/v1/auth/session')
        .catch(() => {})
        .then(() => {
          this.redirectToLogin();
        });
    },
    redirectToLogin() {
      const continueURL = `/${window.location.hash}${window.location.search}`;
      window.location = `/?continue=${encodeURIComponent(continueURL)}`;
    },
    grant(auto = false) {
      this.fetchingGrant = true;

      const data = {
        client_id: this.client_id,
        response_type: this.response_type,
        redirect_uri: this.redirect_uri,
        state: this.state,
        scope: this.scope,
        auto,
      };

      return axios.post('api/v1/auth/oauth2/grant', data)
        .then((resp) => {
          const redirectURL = new URL(this.redirect_uri);
          const params = redirectURL.searchParams;
          if (this.state) params.set('state', this.state);
          params.set('code', resp.data.code);
          redirectURL.search = params.toString();

          setTimeout(() => {
            window.location = redirectURL.toString();
          }, 1500);
        }).catch((err) => {
          this.fetchingGrant = false;
          if (!this.auto) this.error = errorCodes[err.response.data.error];
          throw err;
        });
    },
  },
};
</script>
