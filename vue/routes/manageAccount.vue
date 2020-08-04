<template>
  <div>
    <CenterCard title="Manage Account">
      <AccountInfo />
    </CenterCard>
    <CenterCard title="Audit Log">
      <Audit />
    </CenterCard>
  </div>
</template>

<script>
import CenterCard from '../components/centerCard.vue';
import AccountInfo from '../widgets/accountInfo.vue';
import Audit from '../widgets/audit.vue';

export default {
  components: {
    CenterCard,
    AccountInfo,
    Audit,
  },
};
</script>
