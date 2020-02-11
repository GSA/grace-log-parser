import gzip
import json
import base64
import boto3
import os


def lambda_handler(event, context):

    event_type_list = ["ConsoleLogin"]
    excluded_keys = ["accessKeyId", "principalId"]

    cw_data = event["awslogs"]["data"]

    compressed_payload = base64.b64decode(cw_data)

    uncompressed_payload = gzip.decompress(compressed_payload)

    payload = json.loads(uncompressed_payload)

    print("***PAYLOAD***")
    print(payload)

    log_events = payload["logEvents"]

    for log_event in log_events:
        if "message" in log_event:
            message = json.loads(log_event["message"])
            if "eventName" in message and message["eventName"] in event_type_list:
                ses_message = "This is the header *Test*\n\n"
                ses_message += "EventType: " + message["eventType"] + "\n"
                ses_message += "EventId: " + message["eventID"] + "\n"
                ses_message += "EventTime: " + message["eventTime"] + "\n"
                ses_message += "EventName: " + message["eventName"] + "\n"
                ses_message += "UserAgent: " + message["userAgent"] + "\n"
                ses_message += "AWS Region: " + message["awsRegion"] + "\n"
                ses_message += "SourceIPAddress: " + message["sourceIPAddress"] + "\n"
                if "userIdentity" not in excluded_keys and "userIdentity" in message:
                    ses_message += "\nUserIdentity \n\n"
                    ui_message_builder = []
                    Parse_SES_UserIdentity(
                        message["userIdentity"], ui_message_builder, excluded_keys
                    )
                    ses_message += "\n".join(ui_message_builder)
                push_to_ses(ses_message)


def Parse_SES_UserIdentity(d, ui_message_builder, excluded_keys):
    for k, v in d.items():
        if isinstance(v, dict):
            Parse_SES_UserIdentity(v, ui_message_builder, excluded_keys)
        #   return ses_message
        else:
            if str(k) not in excluded_keys:
                ui_message_builder.append("{0} : {1}".format(k, v))


def push_to_ses(ses_message):
    client = boto3.client("ses")
    response = client.send_email(
        Source=os.environ["FROM_EMAIL"],
        Destination={"ToAddresses": [os.environ["TO_EMAIL"],]},
        Message={
            "Subject": {"Data": os.environ["SUBJECT"]},
            "Body": {"Text": {"Data": ses_message}},
        },
    )
    print("**** SES Publish Response ****")
    print(response)
