<template>
  <CenterCard title="Error">
    <h2 class="title is-2">Page Not Found</h2>
    <p>
      The page you were looking for could not be found. <router-link to="/">Go Home</router-link>
    </p>
  </CenterCard>
</template>

<script>
import CenterCard from '../components/centerCard.vue';

export default {
  components: {
    CenterCard,
  },
};
</script>
