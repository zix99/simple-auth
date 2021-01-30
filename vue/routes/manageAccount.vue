<template>
  <div>
    <CenterCard title="Manage Account">
      <AccountInfo />
    </CenterCard>
    <CenterCard title="OAuth Tokens">
      <AccountOAuth />
    </CenterCard>
    <CenterCard title="Audit Log">
      <AccountAudit />
    </CenterCard>
  </div>
</template>

<script>
import CenterCard from '../components/centerCard.vue';
import AccountInfo from '../widgets/accountInfo.vue';
import AccountAudit from '../widgets/accountAudit.vue';
import AccountOAuth from '../widgets/accountOAuth.vue';

export default {
  components: {
    CenterCard,
    AccountInfo,
    AccountOAuth,
    AccountAudit,
  },
};
</script>
