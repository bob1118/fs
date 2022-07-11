package odbc_cdr

const ODBC_DSN = `<param name="odbc-dsn" value="pgsql://hostaddr=192.168.0.100 dbname=freeswitch user=freeswitch password='freeswitch' options='-c client_min_messages=NOTICE'"/>`
const ODBC_CDR_CONF_XML = `<configuration name="odbc_cdr.conf" description="ODBC CDR Configuration">
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
  <!-- only a-legs will be inserted into this table -->
  <table name="cdr_table_a_leg" log-leg="a-leg">
    <!--  Unique-Id/Variable_uuid -->
    <field name="uuid" chan-var-name="uuid"/>
    <!-- Variable_originated_legs/Variable_originator, aleg/bleg -->
    <field name ="otheruuid" chan-var-name="originated_legs"/>
    <!-- Other-Type: "originatee" aleg-->
    <field name ="othertype" chan-var-name="other_type"/>

    <!-- Variable_channel_name -->
    <field name="ch_name" chan-var-name="channel_name"/>
    <!-- Variable_sofia_profile_name -->
    <field name="ch_profile" chan-var-name="sofia_profile_name"/>
    <!-- Call-Direction/Variable_direction -->
    <field name = "ch_direction" chan-var-name="direction"/>
    <!-- Variable_domain_name -->
    <field name="ch_domain" chan-var-name="domain_name"/>
    <!-- Variable_sip_gateway_name -->
    <field name="ch_gateway" chan-var-name="sip_gateway_name"/>

    <!-- Caller-Caller-ID-Name -->
    <field name ="ch_calleridname" chan-var-name="caller_id_name"/>
    <!-- Caller-Caller-ID-Number -->
    <field name ="ch_calleridnumber" chan-var-name="caller_id_number"/>
    <!-- Caller-Callee-ID-Name -->
    <field name ="ch_calleeidname" chan-var-name="callee_id_name"/>
    <!-- Caller-Callee-ID-Number -->
    <field name ="ch_calleeidnumber" chan-var-name="callee_id_number"/>	
    <!-- Caller-Destination-Number	-->
    <field name ="ch_destination" chan-var-name="destination_number"/>

    <!-- Variable_current_application -->
    <field name ="ch_app" chan-var-name="current_application"/>
    <!-- Variable_current_application_data -->
    <field name ="ch_appdata" chan-var-name="current_application_data"/>
    <!-- Variable_dialstatus -->
    <field name ="ch_dialstatus" chan-var-name="dialstatus"/>
    
    <!-- Variable_hangup_cause = "NORMAL_CLEARING" -->
    <field name ="hangup_cause" chan-var-name="hangup_cause"/>
    <!-- Variable_hangup_cause_q850 = "16" -->
    <field name ="hangup_q850" chan-var-name="hangup_cause_q850"/>
    <!-- Variable_sip_hangup_disposition = "recv_bye";"recv_cancel";"recv_refuse";"send_bye";"send_cancel";"send_refuse"; -->
    <field name ="hangup_disposition" chan-var-name="sip_hangup_disposition"/>
    <!-- Variable_proto_specific_hangup_cause = "sip:200" -->
    <field name ="hangup_protocause" chan-var-name="proto_specific_hangup_cause"/>
    <!-- Variable_sip_hangup_phrase = "OK";"send_bye ok ?"-->
    <field name ="hangup_phrase" chan-var-name="sip_hangup_phrase"/>
    
    <!-- Variable_start_epoch -->
    <field name ="start_epoch" chan-var-name ="start_epoch"/>
    <!-- Variable_answer_epoch -->
    <field name ="answer_epoch" chan-var-name ="answer_epoch"/>
    <!-- Variable_end_epoch -->
    <field name ="end_epoch" chan-var-name ="end_epoch"/>
    <!-- Variable_waitsec -->
    <filed name ="waitsec" chan-var-name ="waitsec"/>
    <!-- Variable_billsec -->
    <field name ="billsec" chan-var-name ="billsec"/>
    <!-- Variable_duration -->
    <field name ="duration" chan-var-name ="duration"/>	
  </table>
  <!-- only b-legs will be inserted into this table -->
  <table name="cdr_table_b_leg" log-leg="b-leg">
    <!--  Unique-Id/Variable_uuid -->
    <field name="uuid" chan-var-name="uuid"/>
    <!-- Variable_originated_legs/Variable_originator, aleg/bleg -->
    <field name ="otheruuid" chan-var-name="originator"/>
    <!-- Other-Type="originator" bleg-->
    <field name ="othertype" chan-var-name="other_type"/>

    <!-- Variable_channel_name -->
    <field name="ch_name" chan-var-name="channel_name"/>
    <!-- Variable_sofia_profile_name -->
    <field name="ch_profile" chan-var-name="sofia_profile_name"/>
    <!-- Call-Direction/Variable_direction -->
    <field name = "ch_direction" chan-var-name="direction"/>
    <!-- Variable_domain_name -->
    <field name="ch_domain" chan-var-name="domain_name"/>
    <!-- Variable_sip_gateway_name -->
    <field name="ch_gateway" chan-var-name="sip_gateway_name"/>

    <!-- Caller-Caller-ID-Name -->
    <field name ="ch_calleridname" chan-var-name="caller_id_name"/>
    <!-- Caller-Caller-ID-Number -->
    <field name ="ch_calleridnumber" chan-var-name="caller_id_number"/>
    <!-- Caller-Callee-ID-Name -->
    <field name ="ch_calleeidname" chan-var-name="callee_id_name"/>
    <!-- Caller-Callee-ID-Number -->
    <field name ="ch_calleeidnumber" chan-var-name="callee_id_number"/>	
    <!-- Caller-Destination-Number	-->
    <field name ="ch_destination" chan-var-name="destination_number"/>

    <!-- Variable_current_application -->
    <field name ="ch_app" chan-var-name="current_application"/>
    <!-- Variable_current_application_data -->
    <field name ="ch_appdata" chan-var-name="current_application_data"/>
    <!-- Variable_dialstatus -->
    <field name ="ch_dialstatus" chan-var-name="dialstatus"/>
    
    <!-- Variable_hangup_cause = "NORMAL_CLEARING" -->
    <field name ="hangup_cause" chan-var-name="hangup_cause"/>
    <!-- Variable_hangup_cause_q850 = "16" -->
    <field name ="hangup_q850" chan-var-name="hangup_cause_q850"/>
    <!-- Variable_sip_hangup_disposition = "recv_bye";"recv_cancel";"recv_refuse";"send_bye";"send_cancel";"send_refuse"; -->
    <field name ="hangup_disposition" chan-var-name="sip_hangup_disposition"/>
    <!-- Variable_proto_specific_hangup_cause = "sip:200" -->
    <field name ="hangup_protocause" chan-var-name="proto_specific_hangup_cause"/>
    <!-- Variable_sip_hangup_phrase = "OK";"send_bye ok ?"-->
    <field name ="hangup_phrase" chan-var-name="sip_hangup_phrase"/>
    
    <!-- Variable_start_epoch -->
    <field name ="start_epoch" chan-var-name ="start_epoch"/>
    <!-- Variable_answer_epoch -->
    <field name ="answer_epoch" chan-var-name ="answer_epoch"/>
    <!-- Variable_end_epoch -->
    <field name ="end_epoch" chan-var-name ="end_epoch"/>
    <!-- Variable_waitsec -->
    <filed name ="waitsec" chan-var-name ="waitsec"/>
    <!-- Variable_billsec -->
    <field name ="billsec" chan-var-name ="billsec"/>
    <!-- Variable_duration -->
    <field name ="duration" chan-var-name ="duration"/>	
  </table>
</tables>
</configuration>`
