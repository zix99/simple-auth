<template>
  <Card title="Manage Account">
    <Login />
  </Card>
</template>

<script>
import Card from './components/card.vue';
import Login from './widgets/login.vue';

export default {
  components: {
    Card,
    Login,
  },
};
</script>
