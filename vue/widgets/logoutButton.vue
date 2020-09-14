<template>
    <button class="button is-primary is-light" @click="logout">Log Out</button>
</template>

<script>
import axios from 'axios';

export default {
  methods: {
    logout() {
      axios.post('api/ui/logout')
        .catch(() => {})
        .then(() => {
          if (this.$router) {
            this.$router.push('/');
          } else {
            window.location = '/';
          }
        });
    },
  },
};
</script>
