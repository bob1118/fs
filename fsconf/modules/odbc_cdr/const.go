package odbc_cdr

const MOD_NAME = `mod_odbc_cdr`
const MOD_CONF_NAME = `odbc_cdr.conf.xml`

const ODBC_DSN = `<param name="odbc-dsn" value="pgsql://hostaddr=192.168.0.100 dbname=freeswitch user=freeswitch password='freeswitch' options='-c client_min_messages=NOTICE'"/>`

const ODBC_CDR_CONF_XML = `
<configuration name="odbc_cdr.conf" description="ODBC CDR Configuration">
<settings>
  <!-- <param name="odbc-dsn" value="database:username:password"/> -->
  <param name="odbc-dsn" value="pgsql://hostaddr=192.168.0.100 dbname=freeswitch user=freeswitch password='freeswitch' options='-c client_min_messages=NOTICE'"/>
  <!-- global value can be "a-leg", "b-leg", "both" (default is "both") -->
  <param name="log-leg" value="both"/>
  <!-- value can be "always", "never", "on-db-fail" -->
  <!-- <param name="write-csv" value="on-db-fail"/> -->
  <!-- location to store csv copy of CDR -->
  <!-- <param name="csv-path" value="/usr/local/freeswitch/log/odbc_cdr"/> -->
  <!-- if "csv-path-on-fail" is set, failed INSERTs will be placed here as CSV files otherwise they will be placed in "csv-path" -->
  <!-- <param name="csv-path-on-fail" value="/usr/local/freeswitch/log/odbc_cdr/failed"/> -->
  <!-- dump SQL statement after leg ends -->
  <param name="debug-sql" value="true"/>
</settings>
<tables>
%s
</tables>
</configuration>
`

const ODBC_CDR_CONF_XML_TABLE_ALEG = `
  <!-- only a-legs will be inserted into this table -->
  <table name="%s" log-leg="a-leg">
    <!-- Variable_uuid -->
    <field name="uuid" chan-var-name="uuid"/>
    <!-- Variable_call_uuid -->
    <field name ="calluuid" chan-var-name="call_uuid"/>
    <!-- Variable_originated_legs/Variable_originator aleg/bleg -->
    <field name ="otheruuid" chan-var-name="originated_legs"/>
    <!-- Other-Type: "originatee"/"originator" aleg/bleg -->
    <!-- <field name ="othertype" chan-var-name="other-type"/> -->
    <!-- Variable_channel_name -->
    <field name="name" chan-var-name="channel_name"/>
    <!-- Variable_sofia_profile_name -->
    <field name="sofiaprofile" chan-var-name="sofia_profile_name"/>
    <!-- Variable_direction -->
    <field name = "direction" chan-var-name="direction"/>
    <!-- Variable_domain_name -->
    <field name="domain" chan-var-name="domain_name"/>
    <!-- Variable_sip_profile_name -->
    <field name="sipprofile" chan-var-name="sip_profile_name"/>
    <!-- Variable_sip_gateway_name -->
    <field name="gateway" chan-var-name="sip_gateway_name"/>

    <!-- Caller-Ani -->
    <field name = "ani" chan-var-name="ani"/>
    <!-- Caller-Destination-Number	-->
    <field name ="destination" chan-var-name="destination_number"/>
    <!-- Caller-Caller-ID-Name -->
    <field name ="calleridname" chan-var-name="caller_id_name"/>
    <!-- Caller-Caller-ID-Number -->
    <field name ="calleridnumber" chan-var-name="caller_id_number"/>
    <!-- Caller-Callee-ID-Name -->
    <field name ="calleeidname" chan-var-name="callee_id_name"/>
    <!-- Caller-Callee-ID-Number -->
    <field name ="calleeidnumber" chan-var-name="callee_id_number"/>	

    <!-- Variable_current_application -->
    <field name ="app" chan-var-name="current_application"/>
    <!-- Variable_current_application_data -->
    <field name ="appdata" chan-var-name="current_application_data"/>
    <!-- Variable_dialstatus -->
    <field name ="dialstatus" chan-var-name="dialstatus"/>
    
    <!-- Variable_hangup_cause = "NORMAL_CLEARING" -->
    <field name ="cause" chan-var-name="hangup_cause"/>
    <!-- Variable_hangup_cause_q850 = "16" -->
    <field name ="q850" chan-var-name="hangup_cause_q850"/>
    <!-- Variable_sip_hangup_disposition = "recv_bye";"recv_cancel";"recv_refuse";"send_bye";"send_cancel";"send_refuse"; -->
    <field name ="disposition" chan-var-name="sip_hangup_disposition"/>
    <!-- Variable_proto_specific_hangup_cause = "sip:200" -->
    <field name ="protocause" chan-var-name="proto_specific_hangup_cause"/>
    <!-- Variable_sip_hangup_phrase = "OK";"send_bye ok ?"-->
    <field name ="phrase" chan-var-name="sip_hangup_phrase"/>
    
    <!-- Variable_start_epoch -->
    <field name ="startepoch" chan-var-name ="start_epoch"/>
    <!-- Variable_answer_epoch -->
    <field name ="answerepoch" chan-var-name ="answer_epoch"/>
    <!-- Variable_end_epoch -->
    <field name ="endepoch" chan-var-name ="end_epoch"/>
    <!-- Variable_waitsec -->
    <filed name ="waitsec" chan-var-name ="waitsec"/>
    <!-- Variable_billsec -->
    <field name ="billsec" chan-var-name ="billsec"/>
    <!-- Variable_duration -->
    <field name ="duration" chan-var-name ="duration"/>	
  </table>
`
const ODBC_CDR_CONF_XML_TABLE_BLEG = `
  <!-- only b-legs will be inserted into this table -->
  <table name="%s" log-leg="b-leg">
    <!-- Variable_uuid -->
    <field name="uuid" chan-var-name="uuid"/>
    <!-- Variable_call_uuid -->
    <field name ="calluuid" chan-var-name="call_uuid"/>
    <!-- Variable_originated_legs/Variable_originator aleg/bleg -->
    <field name ="otheruuid" chan-var-name="originator"/>
    <!-- Other-Type: "originatee"/"originator" aleg/bleg -->
    <!-- <field name ="othertype" chan-var-name="other_type"/> -->
    <!-- Variable_channel_name -->
    <field name="name" chan-var-name="channel_name"/>
    <!-- Variable_sofia_profile_name -->
    <field name="sofiaprofile" chan-var-name="sofia_profile_name"/>
    <!-- Variable_direction -->
    <field name = "direction" chan-var-name="direction"/>
    <!-- Variable_domain_name -->
    <field name="domain" chan-var-name="domain_name"/>
    <!-- Variable_sip_profile_name -->
    <field name="sipprofile" chan-var-name="sip_profile_name"/>
    <!-- Variable_sip_gateway_name -->
    <field name="gateway" chan-var-name="sip_gateway_name"/>

    <!-- Caller-Ani -->
    <field name = "ani" chan-var-name="ani"/>
    <!-- Caller-Destination-Number	-->
    <field name ="destination" chan-var-name="destination_number"/>
    <!-- Caller-Caller-ID-Name -->
    <field name ="calleridname" chan-var-name="caller_id_name"/>
    <!-- Caller-Caller-ID-Number -->
    <field name ="calleridnumber" chan-var-name="caller_id_number"/>
    <!-- Caller-Callee-ID-Name -->
    <field name ="calleeidname" chan-var-name="callee_id_name"/>
    <!-- Caller-Callee-ID-Number -->
    <field name ="calleeidnumber" chan-var-name="callee_id_number"/>	

    <!-- Variable_current_application -->
    <field name ="app" chan-var-name="current_application"/>
    <!-- Variable_current_application_data -->
    <field name ="appdata" chan-var-name="current_application_data"/>
    <!-- Variable_dialstatus -->
    <field name ="dialstatus" chan-var-name="dialstatus"/>
    
    <!-- Variable_hangup_cause = "NORMAL_CLEARING" -->
    <field name ="cause" chan-var-name="hangup_cause"/>
    <!-- Variable_hangup_cause_q850 = "16" -->
    <field name ="q850" chan-var-name="hangup_cause_q850"/>
    <!-- Variable_sip_hangup_disposition = "recv_bye";"recv_cancel";"recv_refuse";"send_bye";"send_cancel";"send_refuse"; -->
    <field name ="disposition" chan-var-name="sip_hangup_disposition"/>
    <!-- Variable_proto_specific_hangup_cause = "sip:200" -->
    <field name ="protocause" chan-var-name="proto_specific_hangup_cause"/>
    <!-- Variable_sip_hangup_phrase = "OK";"send_bye ok ?"-->
    <field name ="phrase" chan-var-name="sip_hangup_phrase"/>
    
    <!-- Variable_start_epoch -->
    <field name ="startepoch" chan-var-name ="start_epoch"/>
    <!-- Variable_answer_epoch -->
    <field name ="answerepoch" chan-var-name ="answer_epoch"/>
    <!-- Variable_end_epoch -->
    <field name ="endepoch" chan-var-name ="end_epoch"/>
    <!-- Variable_waitsec -->
    <filed name ="waitsec" chan-var-name ="waitsec"/>
    <!-- Variable_billsec -->
    <field name ="billsec" chan-var-name ="billsec"/>
    <!-- Variable_duration -->
    <field name ="duration" chan-var-name ="duration"/>	
  </table>
`
const ODBC_CDR_CONF_XML_TABLE_BOTH = `
  <table name="%s" log-leg="both">
    <!--  Variable_uuid -->
    <field name="uuid" chan-var-name="uuid"/>
    <!-- Variable_originator -->
    <field name ="otheruuid" chan-var-name="originator"/>
    <!-- Other-Type: "originatee" aleg-->
    <field name ="othertype" chan-var-name="other_type"/>
  </table>
`
