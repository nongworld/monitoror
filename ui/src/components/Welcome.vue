<template>
  <div class="c-monitoror-welcome" :class="classes">
    <template v-if="shouldShowWelcomeNextStep">
      <h1 class="c-monitoror-welcome--title">Welcome</h1>
      <h2 class="c-monitoror-welcome--subtitle">You're almost there!</h2>

      <div class="c-monitoror-welcome--content">
        <div class="c-monitoror-welcome--next-step-content">
          <div>
            <div class="c-monitoror-welcome--next-step-sup-title">
              1. Next step
            </div>
            <h3 class="c-monitoror-welcome--next-step-title">
              Create your <br>
              UI configuration
            </h3>
            <p>
              Define what you want to display by creating a JSON file using
              tile types.
            </p>
            <p>
              <a href="https://monitoror.com/documentation/#ui-configuration" class="c-monitoror-welcome--button" target="_blank">
                Go to documentation
                <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z" fill="currentColor"></path>
                  <path d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z" fill="currentColor"></path>
                </svg>
              </a>
            </p>
          </div>
          <div>
            <div class="c-monitoror-welcome--next-step-sup-title">
              2. And finally...
            </div>
            <h3 class="c-monitoror-welcome--next-step-title">
              Set your configuration <br>
              as default
            </h3>
            <p>
              Put the path or URL to your JSON file in
              the <code>MO_CONFIG</code> environment variable.
            </p>
            <p>
              <a href="https://monitoror.com/documentation/#core-configuration" class="c-monitoror-welcome--button" target="_blank">
                Read more about Core Configuration
                <svg viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
                  <path d="M31.2 14L18.7 1.6c-1.1-1.1-3-1.1-4.1 0-1.1 1.1-1.1 2.9 0 4.1L25 16 14.6 26.3c-1.1 1.1-1.1 2.9 0 4.1 1.1 1.1 3 1.1 4.1 0L31.2 18c.6-.6.9-1.3.8-2 0-.7-.3-1.5-.8-2z" fill="currentColor"></path>
                  <path d="M2.3 13h22.9c1.3 0 2.3 1 2.3 2.3v1.5c0 1.3-1 2.3-2.3 2.3H2.3c-1.3 0-2.3-1-2.3-2.3v-1.5C0 14.1 1 13 2.3 13z" fill="currentColor"></path>
                </svg>
              </a>
            </p>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <h1 class="c-monitoror-welcome--title">Welcome back</h1>
      <h2 class="c-monitoror-welcome--subtitle">Choose a configuration</h2>

      <div class="c-monitoror-welcome--content c-monitoror-welcome--content__choose-configuration">
        <p v-for="configMetadata in configList" >
          <a :href="configMetadata.uiUrl" class="c-monitoror-welcome--button c-monitoror-welcome--button__card">
            <code>{{configMetadata.name}}</code>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentcolor">
              <path fill-rule="evenodd" d="M11.31 2.525a9.648 9.648 0 011.38 0c.055.004.135.05.162.16l.351 1.45c.153.628.626 1.08 1.173 1.278.205.074.405.157.6.249a1.832 1.832 0 001.733-.074l1.275-.776c.097-.06.186-.036.228 0 .348.302.674.628.976.976.036.042.06.13 0 .228l-.776 1.274a1.832 1.832 0 00-.074 1.734c.092.195.175.395.248.6.198.547.652 1.02 1.278 1.172l1.45.353c.111.026.157.106.161.161a9.653 9.653 0 010 1.38c-.004.055-.05.135-.16.162l-1.45.351a1.833 1.833 0 00-1.278 1.173 6.926 6.926 0 01-.25.6 1.832 1.832 0 00.075 1.733l.776 1.275c.06.097.036.186 0 .228a9.555 9.555 0 01-.976.976c-.042.036-.13.06-.228 0l-1.275-.776a1.832 1.832 0 00-1.733-.074 6.926 6.926 0 01-.6.248 1.833 1.833 0 00-1.172 1.278l-.353 1.45c-.026.111-.106.157-.161.161a9.653 9.653 0 01-1.38 0c-.055-.004-.135-.05-.162-.16l-.351-1.45a1.833 1.833 0 00-1.173-1.278 6.928 6.928 0 01-.6-.25 1.832 1.832 0 00-1.734.075l-1.274.776c-.097.06-.186.036-.228 0a9.56 9.56 0 01-.976-.976c-.036-.042-.06-.13 0-.228l.776-1.275a1.832 1.832 0 00.074-1.733 6.948 6.948 0 01-.249-.6 1.833 1.833 0 00-1.277-1.172l-1.45-.353c-.111-.026-.157-.106-.161-.161a9.648 9.648 0 010-1.38c.004-.055.05-.135.16-.162l1.45-.351a1.833 1.833 0 001.278-1.173 6.95 6.95 0 01.249-.6 1.832 1.832 0 00-.074-1.734l-.776-1.274c-.06-.097-.036-.186 0-.228.302-.348.628-.674.976-.976.042-.036.13-.06.228 0l1.274.776a1.832 1.832 0 001.734.074 6.95 6.95 0 01.6-.249 1.833 1.833 0 001.172-1.277l.353-1.45c.026-.111.106-.157.161-.161zM12 1c-.268 0-.534.01-.797.028-.763.055-1.345.617-1.512 1.304l-.352 1.45c-.02.078-.09.172-.225.22a8.45 8.45 0 00-.728.303c-.13.06-.246.044-.315.002l-1.274-.776c-.604-.368-1.412-.354-1.99.147-.403.348-.78.726-1.129 1.128-.5.579-.515 1.387-.147 1.99l.776 1.275c.042.069.059.185-.002.315a8.45 8.45 0 00-.302.728c-.05.135-.143.206-.221.225l-1.45.352c-.687.167-1.249.749-1.304 1.512a11.149 11.149 0 000 1.594c.055.763.617 1.345 1.304 1.512l1.45.352c.078.02.172.09.22.225.09.248.191.491.303.729.06.129.044.245.002.314l-.776 1.274c-.368.604-.354 1.412.147 1.99.348.403.726.78 1.128 1.129.579.5 1.387.515 1.99.147l1.275-.776c.069-.042.185-.059.315.002.237.112.48.213.728.302.135.05.206.143.225.221l.352 1.45c.167.687.749 1.249 1.512 1.303a11.125 11.125 0 001.594 0c.763-.054 1.345-.616 1.512-1.303l.352-1.45c.02-.078.09-.172.225-.22.248-.09.491-.191.729-.303.129-.06.245-.044.314-.002l1.274.776c.604.368 1.412.354 1.99-.147.403-.348.78-.726 1.129-1.128.5-.579.515-1.387.147-1.99l-.776-1.275c-.042-.069-.059-.185.002-.315.112-.237.213-.48.302-.728.05-.135.143-.206.221-.225l1.45-.352c.687-.167 1.249-.749 1.303-1.512a11.125 11.125 0 000-1.594c-.054-.763-.616-1.345-1.303-1.512l-1.45-.352c-.078-.02-.172-.09-.22-.225a8.469 8.469 0 00-.303-.728c-.06-.13-.044-.246-.002-.315l.776-1.274c.368-.604.354-1.412-.147-1.99-.348-.403-.726-.78-1.128-1.129-.579-.5-1.387-.515-1.99-.147l-1.275.776c-.069.042-.185.059-.315-.002a8.465 8.465 0 00-.728-.302c-.135-.05-.206-.143-.225-.221l-.352-1.45c-.167-.687-.749-1.249-1.512-1.304A11.149 11.149 0 0012 1zm2.5 11a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0zm1.5 0a4 4 0 11-8 0 4 4 0 018 0z"/>
            </svg>
          </a>
        </p>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
  import {Component, Vue} from 'vue-property-decorator'

  import Route from '@/enums/route'
  import ConfigMetadata from '@/interfaces/configMetadata'

  @Component({})
  export default class MonitororWelcome extends Vue {
    /*
     * Data
     */
    private shouldHaveScroll: boolean = false

    /*
     * Computed
     */

    get classes() {
      return {
        'c-monitoror-welcome__no-scroll': !this.shouldHaveScroll,
      }
    }

    get shouldShowWelcomeNextStep(): boolean {
      const isOnWelcomeNextStepRoute = Route.Welcome === this.$store.getters.currentRoute

      return isOnWelcomeNextStepRoute || this.$store.getters.isNewUser
    }

    get configList(): ConfigMetadata[] {
      return this.$store.state.configList
    }

    /*
     * Methods
     */

    private mounted() {
      this.$store.dispatch('fetchConfigList')

      setTimeout(() => {
        this.shouldHaveScroll = true
      }, 1650)
    }
  }
</script>

<style lang="scss">
  .c-monitoror-welcome {
    position: relative;
    padding-top: 260px;
    font-weight: 400;
    font-size: 24px;
    text-align: center;
    animation: welcomeSlideIn 525ms ease 1.15s both;

    &::before {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 460px;
      background: var(--color-dark-background);
      transform: skewY(-7deg) scaleY(2);
      transform-origin: bottom center;
      z-index: -1;
    }

    &__no-scroll .c-monitoror-welcome--next-step-content {
      display: none;
    }

    .c-monitoror-welcome--title {
      display: inline-block;
      margin: 0 auto;
      font-weight: 700;
      font-size: 72px;
      color: transparent;
      background: linear-gradient(to right, var(--color-succeeded), var(--color-action-required));
      -webkit-background-clip: text;
      background-size: 200% 100%;
      animation: welcomeTitleGradient 10s ease 1.7s infinite;
    }

    .c-monitoror-welcome--subtitle {
      font-weight: 300;
      font-size: 32px;
      margin: 0;
    }
  }

  .c-monitoror-welcome--content {
    width: 870px;
    max-width: calc(100% - 50px);
    margin: 60px auto 0;
    padding: 130px 0 0;
    font-size: 20px;
    line-height: 1.6;
    text-align: left;

    &__choose-configuration {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      grid-gap: 20px;
    }
  }

  .c-monitoror-welcome--next-step-sup-title {
    text-transform: uppercase;
    font-size: 18px;
    line-height: 1.2;
    letter-spacing: 1.1px;
    color: var(--color-succeeded);
    transform: translateX(2px);
    margin-bottom: 7px;
    font-family: "JetBrains Mono", monospace;
  }

  .c-monitoror-welcome--next-step-title {
    font-size: 38px;
    line-height: 1.3;
    font-weight: 300;
    margin: 0 0 20px;
    color: #ffffff;
  }

  .c-monitoror-welcome--next-step-content {
    p {
      margin: 10px 0 30px;
    }

    & > * {
      animation: welcomeSlideIn 700ms both;
    }

    & > *:nth-child(2) {
      animation-delay: 400ms;
    }

    @media (max-width: 999px) {
      max-width: 400px;
      margin: 0 auto;

      & > * {
        padding-bottom: 70px;
      }
    }

    @media (min-width: 1000px) {
      display: grid;
      grid-auto-flow: column;
      grid-gap: 10%;
    }
  }

  .c-monitoror-welcome--button {
    display: inline-block;
    text-align: center;
    font-size: 16px;
    font-weight: 600;
    line-height: 1.2;
    letter-spacing: 0.3px;
    text-decoration: none;
    padding: 20px 30px;
    border-radius: 4px;
    border: 1px solid hsla(211, 52%, 88%, 0.15);
    color: var(--color-spindle);
    background: rgba(207, 223, 240, 0.03);
    transition: color 200ms, background 200ms, border-color 200ms;

    &::after {
      content: initial;
    }

    &:hover {
      color: var(--color-background);
      background: var(--color-spindle);
      border-color: var(--color-spindle);
    }

    svg {
      display: inline-block;
      width: 13px;
      margin-left: 10px;
      transform: translateY(1px);
    }

    &__card {
      width: 100%;
      height: 150px;
      text-align: left;
      font-size: 24px;
      font-weight: normal;
      padding: 22px 25px;
      overflow: hidden;

      svg {
        position: absolute;
        right: -20px;
        bottom: -20px;
        width: 100px;
        opacity: 0.1;
      }
    }
  }

  @keyframes welcomeTitleGradient {
    0%,
    100% {
      background-position: 0 50%;
    }
    50% {
      background-position: 100% 50%;
    }
  }

  @keyframes welcomeSlideIn {
    from {
      opacity: 0;
      transform: translateY(10%);
    }
  }
</style>
