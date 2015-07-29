<?php

namespace newsletters\Http\Controllers\Api;

use Illuminate\Http\Request;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Http\Requests\StoreCampaignRequest;
use newsletters\Repositories\CampaignRepository;
use newsletters\Services\CampaignService;

class CampaignController extends Controller
{
    /**
     * @var CampaignService
     */
    private $service;

    public function __construct(CampaignService $service)
    {
        $this->middleware('auth.basic');

        $this->service = $service;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $campaigns = $this->service->findAllCampaigns($request->has('paginate'), 10);

        return response()->json($campaigns, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreCampaignRequest $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function store(StoreCampaignRequest $request)
    {
        $campaign = $this->service->createCampaign($request->all());
        if (isset($campaign)) {
            return response()->json(['status' => 200, 'campaign' => $campaign->id], 200);
        }

        return response()->json(['status' => 412, 'campaign' => ['The specified resource could not be created.']],
            412);
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return \Illuminate\Http\JsonResponse
     */
    public function show($id)
    {
        $campaign = $this->service->findCampaign($id);

        if (isset($campaign)) {
            return response()->json($campaign, 200);
        }

        return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($id)
    {
        if ($this->service->deleteCampaign($id)) {
            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'],
                200);
        }

        return response()->json(['status' => 422, 'campaign' => ['The specified resource could not be deleted.']],
            422);

    }
}
