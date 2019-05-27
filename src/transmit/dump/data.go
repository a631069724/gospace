package dump

import "fmt"

var (
	ORACLE_TAG = "oracle"
)

type ConsumeReq struct {
	BackUrl       string `json:"back_url" oracle:"back_url"`
	AppMerId      string `json:"app_mer_id" oracle:"app_mer_id"`
	TransMerType  string `json:"trans_mer_type" oracle:"trans_mer_type"`
	InstOpId      string `json:"inst_op_id" oracle:"inst_op_id"`
	InstServiceId string `json:"inst_service_id" oracle:"inst_service_id"`
	PayerInfo     string `json:"payer_info" oracle:"payer_info"`
	RiskInfo      string `json:"risk_info" oracle:"risk_info"`
	AccountIdHash string `json:"account_id_hash" oracle:"account_id_hash"`
	DeviceId      string `json:"device_id" oracle:"device_id"`
	CouponInfo    string `json:"coupon_info" oracle:"coupon_info"`
	AreaInfo      string `json:"area_info" oracle:"area_info"`
	OrderId       string `json:"order_id" oracle:"order_id"`
	InstOrderId   string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime       string `json:"txn_time" oracle:"txn_time"`
	TxnAmt        string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode  string `json:"currency_code" oracle:"currency_code"`
	ReqReserved   string `json:"req_reserved" oracle:"req_reserved"`
	Reserved      string `json:"reserved" oracle:"reserved"`
}

type ConsumeResp struct {
	OrderId      string `json:"order_id" oracle:"order_id"`
	InstOrderId  string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime      string `json:"txn_time" oracle:"txn_time"`
	TxnAmt       string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode string `json:"currency_code" oracle:"currency_code"`
	ReqReserved  string `json:"req_reserved" oracle:"req_reserved"`
	Reserved     string `json:"reserved" oracle:"reserved"`
	QueryId      string `json:"query_id" oracle:"query_id"`
}

/*/pay/notify/consume*/
type ConsumeNotifyReq struct {
	MerId              string `json:"mer_id" oracle:"mer_id"`
	TransMerType       string `json:"trans_mer_type" oracle:"trans_mer_type"`
	OrderId            string `json:"order_id" oracle:"order_id"`
	InstOrderId        string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime            string `json:"txn_time" oracle:"txn_time"`
	TxnAmt             string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode       string `json:"currency_code" oracle:"currency_code"`
	ReqReserved        string `json:"req_reserved" oracle:"req_reserved"`
	Reserved           string `json:"reserved" oracle:"reserved"`
	QueryId            string `json:"query_id" oracle:"query_id"`
	ErrCode            string `json:"errcode" oracle:"errcode"`
	ErrMsg             string `json:"errmsg" oracle:"errmsg"`
	SettleAmt          string `json:"settle_amt" oracle:"settle_amt"`
	SettleCurrencyCode string `json:"settle_currency_code" oracle:"settle_currency_code"`
	SettleDate         string `json:"settle_date" oracle:"settle_date"`
	TraceNo            string `json:"trace_no" oracle:"trace_no"`
	TraceTime          string `json:"trace_time" oracle:"trace_time"`
	ExchangeDate       string `json:"exchange_date" oracle:"exchange_date"`
	ExchangeRate       string `json:"exchange_rate" oracle:"exchange_rate"`
	AccNo              string `json:"acc_no" oracle:"acc_no"`
	PayCardIssueName   string `json:"pay_card_issue_name" oracle:"pay_card_issue_name"`
	PayCardType        string `json:"pay_card_type" oracle:"pay_card_type"`
	PayType            string `json:"pay_type" oracle:"pay_type"`
	IssInsCode         string `json:"iss_ins_code" oracle:"iss_ins_code"`
}

type ConsumeNotifyResp struct {
	MerId       string `json:"mer_id" oracle:"mer_id"`
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
}

type ReverseReq struct {
	InstOpId      string `json:"inst_op_id" oracle:"inst_op_id"`
	InstServiceId string `json:"inst_service_id" oracle:"inst_service_id"`
	OrderId       string `json:"order_id" oracle:"order_id"`
	InstOrderId   string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime       string `json:"txn_time" oracle:"txn_time"`
	ReqReserved   string `json:"req_reserved" oracle:"req_reserved"`
	Reserved      string `json:"reserved" oracle:"reserved"`
}

type ReverseResp struct {
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
	ReqReserved string `json:"req_reserved" oracle:"req_reserved"`
	Reserved    string `json:"reserved" oracle:"reserved"`
}

type CancelReq struct {
	BackUrl       string `json:"back_url" oracle:"back_url"`
	InstOpId      string `json:"inst_op_id" oracle:"inst_op_id"`
	InstServiceId string `json:"inst_service_id" oracle:"inst_service_id"`
	OrderId       string `json:"order_id" oracle:"order_id"`
	InstOrderId   string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime       string `json:"txn_time" oracle:"txn_time"`
	OrigQryId     string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId   string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime   string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnAmt        string `json:"txn_amt" oracle:"txn_amt"`
	ReqReserved   string `json:"req_reserved" oracle:"req_reserved"`
	Reserved      string `json:"reserved" oracle:"reserved"`
}

type CancelResp struct {
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
	OrigQryId   string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnAmt      string `json:"txn_amt" oracle:"txn_amt"`
	ReqReserved string `json:"req_reserved" oracle:"req_reserved"`
	Reserved    string `json:"reserved" oracle:"reserved"`
	QueryId     string `json:"query_id" oracle:"query_id"`
}

/*pay/notify/cancel*/
type CancelNotifyReq struct {
	MerId              string `json:"mer_id" oracle:"mer_id"`
	OrderId            string `json:"order_id" oracle:"order_id"`
	InstOrderId        string `json:"inst_order_id" oracle:"inst_order_id"`
	OrigQryId          string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId        string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime        string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnTime            string `json:"txn_time" oracle:"txn_time"`
	TxnAmt             string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode       string `json:"currency_code" oracle:"currency_code"`
	ReqReserved        string `json:"req_reserved" oracle:"req_reserved"`
	Reserved           string `json:"reserved" oracle:"reserved"`
	QueryId            string `json:"query_id" oracle:"query_id"`
	ErrCode            string `json:"errcode" oracle:"errcode"`
	ErrMsg             string `json:"errmsg" oracle:"errmsg"`
	SettleAmt          string `json:"settle_amt" oracle:"settle_amt"`
	SettleCurrencyCode string `json:"settle_currency_code" oracle:"settle_currency_code"`
	SettleDate         string `json:"settle_date" oracle:"settle_date"`
	TraceNo            string `json:"trace_no" oracle:"trace_no"`
	TraceTime          string `json:"trace_time" oracle:"trace_time"`
	ExchangeDate       string `json:"exchange_date" oracle:"exchange_date"`
	ExchangeRate       string `json:"exchange_rate" oracle:"exchange_rate"`
	AccNo              string `json:"acc_no" oracle:"acc_no"`
	IssInsCode         string `json:"iss_ins_code" oracle:"iss_ins_code"`
}

type CancelNotifyResp struct {
	MerId       string `json:"mer_id" oracle:"mer_id"`
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	OrigQryId   string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
}

type ReturnReq struct {
	BackUrl       string `json:"back_url" oracle:"back_url"`
	InstOpId      string `json:"inst_op_id" oracle:"inst_op_id"`
	InstServiceId string `json:"inst_service_id" oracle:"inst_service_id"`
	OrderId       string `json:"order_id" oracle:"order_id"`
	InstOrderId   string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime       string `json:"txn_time" oracle:"txn_time"`
	OrigQryId     string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId   string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime   string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnAmt        string `json:"txn_amt" oracle:"txn_amt"`
	ReqReserved   string `json:"req_reserved" oracle:"req_reserved"`
	Reserved      string `json:"reserved" oracle:"reserved"`
}

type ReturnResp struct {
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
	OrigQryId   string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnAmt      string `json:"txn_amt" oracle:"txn_amt"`
	ReqReserved string `json:"req_reserved" oracle:"req_reserved"`
	Reserved    string `json:"reserved" oracle:"reserved"`
	QueryId     string `json:"query_id" oracle:"query_id"`
}

/*pay/notify/return*/
type ReturnNotifyReq struct {
	MerId              string `json:"mer_id" oracle:"mer_id"`
	OrderId            string `json:"order_id" oracle:"order_id"`
	InstOrderId        string `json:"inst_order_id" oracle:"inst_order_id"`
	OrigQryId          string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId        string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime        string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnTime            string `json:"txn_time" oracle:"txn_time"`
	TxnAmt             string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode       string `json:"currency_code" oracle:"currency_code"`
	ReqReserved        string `json:"req_reserved" oracle:"req_reserved"`
	Reserved           string `json:"reserved" oracle:"reserved"`
	QueryId            string `json:"query_id" oracle:"query_id"`
	ErrCode            string `json:"errcode" oracle:"errcode"`
	ErrMsg             string `json:"errmsg" oracle:"errmsg"`
	SettleAmt          string `json:"settle_amt" oracle:"settle_amt"`
	SettleCurrencyCode string `json:"settle_currency_code" oracle:"settle_currency_code"`
	SettleDate         string `json:"settle_date" oracle:"settle_date"`
	TraceNo            string `json:"trace_no" oracle:"trace_no"`
	TraceTime          string `json:"trace_time" oracle:"trace_time"`
	ExchangeDate       string `json:"exchange_date" oracle:"exchange_date"`
	ExchangeRate       string `json:"exchange_rate" oracle:"exchange_rate"`
	AccNo              string `json:"acc_no" oracle:"acc_no"`
	IssInsCode         string `json:"iss_ins_code" oracle:"iss_ins_code"`
}

type ReturnNotifyResp struct {
	MerId       string `json:"mer_id" oracle:"mer_id"`
	OrderId     string `json:"order_id" oracle:"order_id"`
	InstOrderId string `json:"inst_order_id" oracle:"inst_order_id"`
	OrigQryId   string `json:"orig_qry_id" oracle:"orig_qry_id"`
	OrigOrderId string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime string `json:"orig_txn_time" oracle:"orig_txn_time"`
	TxnTime     string `json:"txn_time" oracle:"txn_time"`
}

type QueryReq struct {
	QryType  string `json:"qry_type" oracle:"qry_type"`
	QueryId  string `json:"query_id" oracle:"query_id"`
	OrderId  string `json:"order_id" oracle:"order_id"`
	TxnTime  string `json:"txn_time" oracle:"txn_time"`
	Reserved string `json:"reserved" oracle:"reserved"`
}

type QueryResp struct {
	MerId              string `json:"mer_id" oracle:"mer_id"`
	TransMerType       string `json:"trans_mer_type" order:"trans_mer_type"`
	OrderId            string `json:"order_id" oracle:"order_id"`
	InstOrderId        string `json:"inst_order_id" oracle:"inst_order_id"`
	TxnTime            string `json:"txn_time" oracle:"txn_time"`
	TxnAmt             string `json:"txn_amt" oracle:"txn_amt"`
	CurrencyCode       string `json:"currency_code" oracle:"currency_code"`
	ReqReserved        string `json:"req_reserved" oracle:"req_reserved"`
	Reserved           string `json:"reserved" oracle:"reserved"`
	QueryId            string `json:"query_id" oracle:"query_id"`
	OrigOrderId        string `json:"orig_order_id" oracle:"orig_order_id"`
	OrigTxnTime        string `json:"orig_txn_time" oracle:"orig_txn_time"`
	SettleAmt          string `json:"settle_amt" oracle:"settle_amt"`
	SettleCurrencyCode string `json:"settle_currency_code" oracle:"settle_currency_code"`
	SettleDate         string `json:"settle_date" oracle:"settle_date"`
	TraceNo            string `json:"trace_no" oracle:"trace_no"`
	TraceTime          string `json:"trace_time" oracle:"trace_time"`
	ExchangeDate       string `json:"exchange_date" oracle:"exchange_date"`
	ExchangeRate       string `json:"exchange_rate" oracle:"exchange_rate"`
	OrigErrCode        string `json:"orig_errcode" oracle:"orig_errcode"`
	OrigErrMsg         string `json:"orig_errmsg" oracle:"orig_errmsg"`
	AccNo              string `json:"acc_no" oracle:"acc_no"`
	PayCardType        string `json:"pay_card_type" oracle:"pay_card_type"`
	PayType            string `json:"pay_type" oracle:"pay_type"`
	IssInsCode         string `json:"iss_ins_code" oracle:"iss_ins_code"`
}

func OracleTableInit(storer *OraclStorer) {
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_consume_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_consume_req (
   back_url varchar(256),
   app_mer_id varchar(15),
   trans_mer_type varchar(2),
   inst_op_id varchar(16),
   inst_service_id varchar(2),
   payer_info varchar(256),
   risk_info varchar(256),
   account_id_hash varchar(64),
   device_id varchar(64),
   coupon_info varchar(256),
   area_info varchar(7),
   order_id varchar(32),
   inst_order_id varchar(32),
   txn_time varchar(14),
   txn_amt varchar(12),
   currency_code varchar(3),
   req_reserved varchar(512),
   reserved varchar(512) 
  )';
  alter table PAY_CONSUME_REQ
  add constraint PK_CONSUME_REQ primary key (APP_MER_ID, TRANS_MER_TYPE, INST_OP_ID, INST_SERVICE_ID, DEVICE_ID, ORDER_ID, TXN_TIME, INST_ORDER_ID)
  using index ;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_consume_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_consume_resp (
   order_id varchar(32),
   inst_order_id varchar(32),
   txn_time varchar(14),
   txn_amt varchar(12),
   currency_code varchar(3),
   req_reserved varchar(512),
   reserved varchar(512),
   query_id varchar(21)
  )';
  alter table PAY_CONSUME_RESP
  add constraint PK_CONSUME_RESP primary key (ORDER_ID, INST_ORDER_ID, TXN_TIME, QUERY_ID)
  using index ;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_consume_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_consume_req (
   mer_id varchar2(15),
   trans_mer_type varchar2(2),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   txn_amt varchar2(12),
   currency_code varchar2(3),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21),
   errcode varchar2(4),
   errmsg varchar2(1024),
   settle_amt varchar2(12),
   settle_currency_code varchar2(3),
   settle_date varchar2(4),
   trace_no varchar2(6),
   trace_time varchar2(10),
   exchange_date varchar2(4),
   exchange_rate varchar2(8),
   acc_no varchar2(1024),
   pay_card_issue_name varchar2(64),
   pay_card_type varchar2(2),
   pay_type varchar2(4),
   iss_ins_code varchar2(11)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_consume_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_consume_resp (
   mer_id varchar2(15),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_reverse_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_reverse_req (
   inst_op_id varchar2(16),
   inst_service_id varchar2(2),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   req_reserved varchar2(512),
   reserved varchar2(512)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_reverse_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_reverse_resp (
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   req_reserved varchar2(512),
   reserved varchar2(512)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_cancel_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_cancel_req (
   back_url varchar2(256),
   inst_op_id varchar2(16),
   inst_service_id varchar2(2),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_amt varchar2(12),
   req_reserved varchar2(512),
   reserved varchar2(512)
  )';
  alter table PAY_CANCEL_REQ
  add constraint PK_CANCEL_REQ primary key (INST_OP_ID, INST_SERVICE_ID, ORDER_ID, TXN_TIME, INST_ORDER_ID)
  using index;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_cancel_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_cancel_resp (
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_amt varchar2(12),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21)
  )';
  alter table PAY_CANCEL_RESP
  add constraint PK_CANCEL_RESP primary key (ORDER_ID, INST_ORDER_ID, TXN_TIME, ORIG_QRY_ID, ORIG_ORDER_ID)
  using index ;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_cancel_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_cancel_req (
   mer_id varchar2(15),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_time varchar2(14),
   txn_amt varchar2(12),
   currency_code varchar2(3),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21),
   errcode varchar2(4),
   errmsg varchar2(1024),
   settle_amt varchar2(12),
   settle_currency_code varchar2(3),
   settle_date varchar2(4),
   trace_no varchar2(6),
   trace_time varchar2(10),
   exchange_date varchar2(4),
   exchange_rate varchar2(8),
   acc_no varchar2(1024),
   iss_ins_code varchar2(11)
  )';
  alter table PAY_NOTIFY_CANCEL_REQ
  add constraint PK_NOTIFY_CONCEL_REQ primary key (MER_ID, ORDER_ID, INST_ORDER_ID, ORIG_QRY_ID, ORIG_TXN_TIME, TXN_TIME, ORIG_ORDER_ID, QUERY_ID, TRACE_NO, TRACE_TIME)
  using index;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_cancel_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_cancel_resp (
   mer_id varchar2(15),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_time varchar2(14)
  )';
  alter table PAY_NOTIFY_CANCEL_RESP
  add constraint PK_NOTIFY_CANCEL_RESP primary key (MER_ID, ORDER_ID, INST_ORDER_ID, ORIG_QRY_ID, ORIG_ORDER_ID, TXN_TIME)
  using index;
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_return_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_return_req (
   back_url varchar2(256),
   inst_op_id varchar2(16),
   inst_service_id varchar2(2),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_amt varchar2(12),
   req_reserved varchar2(512),
   reserved  varchar2(512)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_return_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_return_resp (
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_amt varchar2(12),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_return_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_return_req (
   mer_id varchar2(15),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_time varchar2(14),
   txn_amt varchar2(12),
   currency_code varchar2(3),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21),
   errcode varchar2(4),
   errmsg varchar2(1024),
   settle_amt varchar2(12),
   settle_currency_code varchar2(3),
   settle_date varchar2(4),
   trace_no varchar2(6),
   trace_time varchar2(10),
   exchange_date varchar2(4),
   exchange_rate varchar2(8),
   iss_ins_code varchar2(11)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_return_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_return_resp (
   mer_id varchar2(15),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   orig_qry_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   txn_time varchar2(14)
  )';
   END IF;
 END; `)
	storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_query_req');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_query_req (
   qry_type varchar2(2),
   query_id varchar2(21),
   order_id varchar2(32),
   txn_time varchar2(14),
   reserved varchar2(512)
  )';
   END IF;
 END; `)
	err := storer.RunSql(`DECLARE
   TOTAL INT := 0;
 BEGIN
   SELECT COUNT(1)
     INTO TOTAL
     FROM USER_TABLES A
    WHERE A.TABLE_NAME = upper('pay_notify_query_resp');
   IF TOTAL = 0 THEN
     EXECUTE IMMEDIATE '
  create table pay_notify_query_resp (
   mer_id varchar2(15),
   trans_mer_type varchar2(2),
   order_id varchar2(32),
   inst_order_id varchar2(32),
   txn_time varchar2(14),
   txn_amt varchar2(12),
   currency_code varchar2(3),
   req_reserved varchar2(512),
   reserved varchar2(512),
   query_id varchar2(21),
   orig_order_id varchar2(32),
   orig_txn_time varchar2(14),
   orig_errcode varchar2(4),
   orig_errmsg varchar2(1024),
   settle_amt varchar2(12),
   settle_currency_code varchar2(3),
   settle_date varchar2(4),
   trace_no varchar2(6),
   trace_time varchar2(10),
   exchange_date varchar2(4),
   exchange_rate varchar2(8),
   acc_no varchar2(1024),
   pay_card_type varchar2(2),
   pay_type varchar2(4),
   iss_ins_code varchar2(11)
  )';
   END IF;
 END; `)
	if err != nil {
		fmt.Println("create table err:", err.Error())
	}

}
