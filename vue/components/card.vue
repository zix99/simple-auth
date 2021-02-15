<template>
    <div class="card my-3">
        <header class="card-header" v-if="title">
            <p class="card-header-title">{{title}}</p>
        </header>
        <div class="card-content">
            <div class="content">
                <slot></slot>
            </div>
        </div>
    </div>
</template>

<script>
export default {
  props: {
    title: null,
  },
};
</script>
