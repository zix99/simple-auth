<template>
  <div>
    <CenterCard title="Manage Account">
      <LoadingBanner :promise="loadingPromise">
        Fetching account details...
      </LoadingBanner>

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
    <CenterCard title="Audit Log">
      <Audit />
    </CenterCard>
  </div>
</template>

<script>
import axios from 'axios';
import Card from '../components/card.vue';
import CenterCard from '../components/centerCard.vue';
import LogoutButton from '../widgets/logoutButton.vue';
import Audit from '../widgets/audit.vue';
import LoadingBanner from '../components/loadingBanner.vue';

export default {
  components: {
    Card,
    CenterCard,
    LogoutButton,
    LoadingBanner,
    Audit,
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
      });
  },
};
</script>
