package switch_modules

const PRE_LOAD_MODULES_CONF = `
<configuration name="pre_load_modules.conf" description="Modules">
  <modules>
    <!-- Databases -->
    <!-- <load module="mod_mariadb"/> -->
    <load module="mod_pgsql"/>
  </modules>
</configuration>
`
const POST_LOAD_MODULES_CONF = `
<configuration name="post_load_modules.conf" description="Modules">
  <modules>
  </modules>
</configuration>
`
