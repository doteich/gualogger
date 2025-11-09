<script setup lang="ts">

</script>

<template>
  <div class="logo-container ">
    <p>
      <span>G</span>eist
    </p>

  </div>
</template>

<style scoped>
.logo-container {
  display: flex;
  font-family: "Zen Dots";
  justify-content: center;
  align-items: center;
  height: 100vh;
  font-size: 1.5rem;
}



.logo-container>p>span {
  font-size: 8rem;
  display: inline-block;
  transform: rotate(45deg);
  text-shadow: -1px -5px 5px rgb(255, 0, 119), -1px -10px 5px rgb(0, 140, 255);
}
</style>
