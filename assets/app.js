$(function(){

	var imageURL = Arg.get("imageURL")
	if (imageURL && imageURL.length > 0) {
		findSimilar(imageURL)
	} else {
		showSamples()
	}

	function render(response) {
		console.info(response)
		for (var i in response.items) {
			var item = response.items[i]
			console.info(item)

			var niceConfidence = "";
			if (item.confidence) {
				niceConfidence = " ("+Math.floor(item.confidence * 100)+"% confidence)"
			}
			var tooltip = item.title
			if (tooltip.length > 60) {
				tooltip = tooltip.substring(0, 60)+"..."
			}
			tooltip += niceConfidence

			$("#results").append(
				$("<a>", {
					'href':"/?imageURL="+encodeURIComponent(item.url),
					'class':'round image',
					'data-tooltip': tooltip
				}).css({
					'background-image': 'url('+item.url+')'
				})
			)

			// $("#results").append(
			// 	$("<div>", {class:"card"}).append(
			// 		$("<a>", {
			// 			class:"image",
			// 			href:"/?imageURL="+encodeURIComponent(item.url),
			// 			'data-tooltip': item.title+' '+niceConfidence
			// 		}).append(
			// 			$("<img>", {src: item.url})
			// 		)
			// 		// $("<div>", {class:"content"}).append(
			// 		// 	//$("<div>", {class:""}).text(item.title),
			// 		// 	$("<span>", {class:"meta",'data-tooltip':item.confidence}).text(niceConfidence)
			// 		// )
			// 	)
			// )
		}
	}

	function findSimilar(imageURL) {
		$('input[name="imageURL"]').val(imageURL)
		$('body').css({
			backgroundImage: 'url('+imageURL+')'
		})
		$('#results').empty()
		$.ajax({
			url: '/api/similar-images',
			data: {
				'url': imageURL
			},
			success: render,
			error: function(){
				console.warn(arguments)
				alert("oops, something went wrong (check the console)")
			}
		})
	}

	function showSamples() {
		$('#results').empty()
		$.ajax({
			url: '/api/random-images',
			success: render,
			error: function(){
				console.warn(arguments)
				alert("oops, something went wrong (check the console)")
			}
		})
	}

})
