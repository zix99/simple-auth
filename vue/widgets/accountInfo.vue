<template>
  <div>
    <LoadingBanner :promise="loadingPromise">
      Fetching account details...
    </LoadingBanner>
    <div v-if="account">
      <h2 class="title is-3">{{account.email}}</h2>
      <Card v-if="account.auth.simple" title="Simple Auth">
        <table class="table">
          <tbody>
            <tr>
              <th>Enabled</th><td>{{account.auth.simple ? 'Yes' : 'No'}}</td>
            </tr>
            <tr>
              <th>Username</th><td>{{account.auth.simple.username}}</td>
            </tr>
            <tr>
              <th>Password</th><td><button class="button is-warning is-light" @click="$refs.modalPass.open()">Change Password</button></td>
            </tr>
          </tbody>
        </table>
      </Card>
      <Card v-for="oidc in account.auth.oidc" :key="oidc.provider" :title="`${oidc.name} Auth`">
        <table class="table">
          <tbody>
            <tr>
              <th>Provider</th><td><i class="fab" :class="oidc.icon"></i> {{oidc.name}}</td>
            </tr>
            <tr>
              <th>Subject</th><td>{{oidc.subject}}</td>
            </tr>
          </tbody>
        </table>
      </Card>
      <div class="box has-text-centered">
        <LogoutButton />
      </div>
    </div>

    <Modal ref="modalPass" title="Update Password">
      <ChangePassword @submitted="$refs.modalPass.close()" />
    </Modal>
  </div>
</template>

<script>
import axios from 'axios';
import Card from '../components/card.vue';
import LoadingBanner from '../components/loadingBanner.vue';
import Modal from '../components/modal.vue';
import LogoutButton from './logoutButton.vue';
import ChangePassword from './changePassword.vue';

export default {
  components: {
    Card,
    Modal,
    LoadingBanner,
    LogoutButton,
    ChangePassword,
  },
  data() {
    return {
      account: null,
      loadingPromise: null,
    };
  },
  created() {
    this.loadingPromise = axios.get('/api/ui/account')
      .then((resp) => {
        this.account = resp.data;
      }).catch(() => {
        this.$router.push('/');
      });
  },
};
</script>
