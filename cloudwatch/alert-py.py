import json
import urllib3
import os
import logging

# 配置日志记录
logger = logging.getLogger()
logger.setLevel(logging.INFO)

http = urllib3.PoolManager()
TELEGRAM_TOKEN = ''
TELEGRAM_CHAT_ID = ''

def lambda_handler(event, context):
    try:
        # 确保 'Records' 和 'Sns' 字段存在
        if 'Records' in event and event['Records'][0].get('Sns'):
            alarm_message = event['Records'][0]['Sns'].get('Message', '{}')
            logger.info(f"Received alarm message: {alarm_message}")
            alarm_data = json.loads(alarm_message)
        else:
            logger.error("Invalid event format: Missing 'Records' or 'Sns'")
            return {
                'statusCode': 400,
                'body': json.dumps('Invalid event format')
            }

        # 提取告警相关信息
        alarm_name = alarm_data.get('AlarmName', 'N/A')
        state = alarm_data.get('NewStateValue', 'N/A')
        reason = alarm_data.get('NewStateReason', 'N/A')
        region = alarm_data.get('Region', 'N/A')




        # 自定义告警消息格式
        custom_message = f"""
🚨 *CloudWatch 告警* 🚨

*告警信息*: {alarm_name}
*当前状态*: {state}
*详细信息*: {reason}
*所在地区*: {region}
"""

        logger.info(f"Sending custom message to Telegram: {custom_message}")
        send_to_telegram(custom_message)

        return {
            'statusCode': 200,
            'body': json.dumps('Success')
        }

    except Exception as e:
        logger.error(f"Error processing CloudWatch alarm: {str(e)}", exc_info=True)
        return {
            'statusCode': 500,
            'body': json.dumps(f'Error: {str(e)}')
        }


def send_to_telegram(message):
    """发送告警消息到 Telegram"""
    url = f"https://api.telegram.org/bot{TELEGRAM_TOKEN}/sendMessage"
    payload = {
        'chat_id': TELEGRAM_CHAT_ID,
        'text': message,
        'parse_mode': 'Markdown'  # 使用 Markdown 格式化消息
    }

    try:
        response = http.request(
            'POST',
            url,
            headers={'Content-Type': 'application/json'},
            body=json.dumps(payload)
        )

        if response.status != 200:
            logger.error(f"Failed to send message to Telegram: {response.status}")
        else:
            logger.info(f"Message sent to Telegram successfully with status: {response.status}")
    except Exception as e:
        logger.error(f"Exception when sending message to Telegram: {str(e)}", exc_info=True)

