<template>
  <CenterCard title="Manage Account">
    <LoadingBanner v-if="loading">
      Fetching account details...
    </LoadingBanner>

    <div v-if="error">
      <article class="message is-danger">
        <div class="message-body">{{error}}</div>
      </article>
      <router-link to="/">Go Home</router-link>
    </div>

    <div v-if="account">
      <h2 class="title is-3">{{account.email}}</h2>
      <Card v-if="account.auth.simple" title="Simple Auth">
        <table class="table">
          <tbody>
            <tr>
              <th>Enabled</th><td>Yes</td>
            </tr>
            <tr>
              <th>Username</th><td>{{account.auth.simple.username}}</td>
            </tr>
          </tbody>
        </table>
      </Card>

      <div class="box has-text-centered">
        <LogoutButton />
      </div>
    </div>
  </CenterCard>
</template>

<script>
import axios from 'axios';
import Card from '../components/card.vue';
import CenterCard from '../components/centerCard.vue';
import LogoutButton from '../widgets/logoutButton.vue';
import LoadingBanner from '../components/loadingBanner.vue';

export default {
  components: {
    Card,
    CenterCard,
    LogoutButton,
    LoadingBanner,
  },
  data() {
    return {
      account: null,
      error: null,
      loading: false,
    };
  },
  created() {
    this.loading = true;
    axios.get('/api/ui/account')
      .then((resp) => {
        this.account = resp.data;
      }).catch((err) => {
        this.error = err.message;
      }).then(() => {
        this.loading = false;
      });
  },
};
</script>
