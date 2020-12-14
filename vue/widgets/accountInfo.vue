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
              <th>Password</th><td><button class="button is-warning is-light" @click="$refs.modalPass.open()">Change</button></td>
            </tr>
            <tr v-if="account.auth.simple.twofactorallowed">
              <th>Two Factor</th>
              <td v-if="!account.auth.simple.twofactor">
                <button class="button is-secondary is-light" @click="$refs.modalTFA.open()">Activate</button>
              </td>
              <td v-else>
                <button class="button is-danger is-light" @click="$refs.deactivateTFA.open()">Deactivate</button>
              </td>
            </tr>
          </tbody>
        </table>
      </Card>
      <Card v-for="oidc in account.auth.oidc" :key="oidc.provider" :title="`${oidc.name} Auth`">
        <table class="table">
          <tbody>
            <tr>
              <th>Provider</th><td><fa-icon :icon="['fab', oidc.icon]" /> {{oidc.name}}</td>
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

    <Modal ref="modalTFA" title="Activate Two Factor">
      <TwoFactorSetup @submitted="$refs.modalTFA.close(); refresh()" />
    </Modal>
    <Modal ref="deactivateTFA" title="Deactivate Two Factor">
      <TwoFactorDeactivate @submitted="$refs.deactivateTFA.close(); refresh()" />
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
import TwoFactorSetup from './twoFactorSetup.vue';
import TwoFactorDeactivate from './twoFactorDeactivate.vue';

export default {
  components: {
    Card,
    Modal,
    LoadingBanner,
    LogoutButton,
    ChangePassword,
    TwoFactorSetup,
    TwoFactorDeactivate,
  },
  data() {
    return {
      account: null,
      loadingPromise: null,
    };
  },
  created() {
    this.refresh();
  },
  methods: {
    refresh() {
      this.loadingPromise = axios.get('api/v1/account')
        .then((resp) => {
          this.account = resp.data;
        }).catch(() => {
          this.$router.push('/');
        });
    },
  },
};
</script>
