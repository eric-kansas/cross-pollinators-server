# Copyright 2016, Beeswax.IO Inc.

'''
Beeswax Win Log Requester.

Generate Beeswax win logs (Impression, Click, Activity) http requests to a designated endpoint
with body generated from the specified input file. Output the http responses to a file
if specified.
'''

from argparse import ArgumentParser, ArgumentDefaultsHelpFormatter
import logging
import sys

from google.protobuf import text_format
import requests

from beeswax.logs.ad_log_pb2 import AdLogMessage


logger = logging.getLogger('beeswax.win_log_requester')


ITEM_DELIMITER = '==='
_HTTP_TIMEOUT_S = 0.1


_LOG_LEVEL_NAMES = ['error', 'info', 'debug']


def parse_command_line():
    parser = ArgumentParser('Beeswax Win Log Requester', formatter_class=ArgumentDefaultsHelpFormatter)
    parser.add_argument('path_to_requests_file')
    parser.add_argument('log_endpoint', help='Log endpoint.')
    parser.add_argument('--path-to-responses-file',
                        help='Path to log responses output file. '
                        'will skip writing response if not present')
    parser.add_argument('--header-secret', help='Authentication secret in log header')
    parser.add_argument('--log-level', default='info', choices=_LOG_LEVEL_NAMES)
    parser.add_argument('--json-format', action='store_true', help='assume input is json (otherwise assume protobuf debug strings)')
    return parser.parse_args()


def _request_text_generator(requests_input_file):

    '''Generate win log requests in ASCII from input file'''

    while True:
        request_text_buffer = []
        for line in requests_input_file:
            if line.startswith(ITEM_DELIMITER):
                break
            request_text_buffer.append(line)
        if not request_text_buffer:
            break
        yield ''.join(request_text_buffer)


def _write_response(responses_output_file, response_message):

    '''Write response_info to responses_output_file if responses_output_file is not None'''

    if responses_output_file:
        responses_output_file.write('{}\n'.format(response_message))
        responses_output_file.write('{}\n'.format(ITEM_DELIMITER))


def _get_response_message(response):

    '''Return response info from http response'''

    response_msg = '''<Status> [{}] {}
<Response headers>
{}
'''.format(response.status_code, response.reason, response.headers)

    return response_msg


def main():
    opts = parse_command_line()
    logger.setLevel(logging._levelNames[opts.log_level.upper()])
    logger.addHandler(logging.StreamHandler())

    logger.info('Endpoint: {}'.format(opts.log_endpoint))

    if opts.json_format:
        headers = {
            'Content-type': 'application/json',
        }
    else:
        headers = {
            'Content-type': 'application/x-protobuf',
        }

    if opts.header_secret:
        headers['beeswax-auth-secret'] = opts.header_secret

    try:
        input_request_file = open(opts.path_to_requests_file, 'rb')
    except (IOError, OSError) as exc:
        logger.error('Could not open impression log requests input file: {}'.format(exc))
        return -1

    output_file = None
    if opts.path_to_responses_file:
        try:
            output_file = open(opts.path_to_responses_file, 'wb')
        except (IOError, OSError) as exc:
            logger.error('Could not open win log responses output file: {}'.format(exc))
            return -1

    try:
        session = requests.Session()
        session.headers.update(headers)

        success_count = 0
        failure_count = 0
        warning_count = 0
        request_data = None

        for request_text in _request_text_generator(input_request_file):
            try:
                # using json format
                if opts.json_format:
                    request_data = request_text
                # using protobuf
                else:
                    request_proto = AdLogMessage()
                    text_format.Parse(request_text, request_proto)
                    request_data = request_proto.SerializeToString()

            except Exception as exc:
                msg = 'Could not parse win log: {} \n The request is: {}'.format(exc, request_text)
                logger.error(msg)
                # Intentionally write errors into output file so that (1) responses (errors) will
                # be aligned with requests and (2) user can do analysis in the output file.
                _write_response(output_file, msg)
                failure_count += 1
                continue

            try:
                logger.debug('Sending win log: {}'.format(request_data))
                response = session.post(opts.log_endpoint,
                                        data=request_data,
                                        timeout=_HTTP_TIMEOUT_S)

            except Exception as exc:
                msg = 'Error in sending http request: {}'.format(exc)
                logger.error(msg)
                # Intentionally write errors into output file.
                _write_response(output_file, msg)
                failure_count += 1
                continue

            if response.status_code not in (200, 204):
                msg = 'Not expected response status code: {}'.format(response.status_code)
                # Intentionally write warnings into output file.
                _write_response(output_file, msg)
                warning_count += 1
                continue

            _write_response(output_file, _get_response_message(response))
            success_count += 1
            logger.debug('Successfully processed request: {}'.format(request_text))

        input_request_file.close()
    finally:
        if output_file:
            output_file.close()

        logger.info('Finished processing all requests. Success count: {}, failure count: {}, warning count: {}'
                    .format(success_count, failure_count, warning_count))

    return 0


if __name__ == '__main__':
    sys.exit(main())
